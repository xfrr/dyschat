package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/events"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
	"github.com/xfrr/dyschat/messages/commands"
	"github.com/xfrr/dyschat/messages/telemetry"
	"github.com/xfrr/dyschat/pkg/env"
	"github.com/xfrr/dyschat/pkg/log"
	"github.com/xfrr/dyschat/pkg/telemetry/jaeger"

	icommands "github.com/xfrr/dyschat/internal/commands"
	inats "github.com/xfrr/dyschat/internal/pubsub/nats"
	mnats "github.com/xfrr/dyschat/messages/nats"
	iredis "github.com/xfrr/dyschat/messages/redis"
	tpkg "github.com/xfrr/dyschat/pkg/telemetry"
	tprom "github.com/xfrr/dyschat/pkg/telemetry/prometheus"
	xprom "go.opentelemetry.io/otel/exporters/prometheus"
)

const (
	serviceName          string = "messages-writer"
	messagesNatsStream   string = "messages"
	roomEventsNatsStream string = "room_events"
)

var (
	svcName = "messages-writer"

	// common
	jaegerURL = env.Get("DYCHAT_JAEGER_URL", "http://jaeger:14268/api/traces")
	natsURL   = env.Get("DYCHAT_NATS_URL", "nats://0.0.0.0:4222")
	redisAddr = env.Get("DYCHAT_REDIS_ADDR", "redis:6379")
	redisPass = env.Get("DYCHAT_REDIS_PASS", "")
	redisDB   = env.Get("DYCHAT_REDIS_DATABASE", "0")

	// messages-writer
	envMode  = env.Get("DYCHAT_MSG_WRITER_ENV_MODE", "development")
	logLevel = env.Get("DYCHAT_MSG_WRITER_LOG_LEVEL", "debug")
)

func main() {
	logger := log.NewZeroLogger(log.ParseLogLevel(logLevel))

	idp := idn.NewNanoIDProvider()

	nc, err := inats.Connect(natsURL)
	if err != nil {
		panic(err)
	}
	defer nc.Drain()

	streams := inats.NewStreams(
		nats.StreamConfig{
			Name: messagesNatsStream,
			Subjects: []string{
				"rooms.*.messages.>",
			},
			Retention: nats.InterestPolicy,
		},
		nats.StreamConfig{
			Name: roomEventsNatsStream,
			Subjects: []string{
				"rooms.*.events.>",
				"rooms.*.members.>",
			},
			Retention: nats.InterestPolicy,
		},
	)
	jsc, err := inats.NewJetStreamContext(nc, streams)
	if err != nil {
		panic(err)
	}

	db, _ := strconv.Atoi(redisDB)
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       db,
	})

	tracer, err := jaeger.NewTracerProvider(jaegerURL, svcName, envMode)
	if err != nil {
		panic(err)
	}

	exporter, err := xprom.New()
	if err != nil {
		panic(err)
	}

	meter := tprom.NewMeterProvider(svcName, exporter)
	go serveMetrics()

	stge := telemetry.NewStorageTelemetryMiddleware(
		tracer,
		meter,
		iredis.NewStorage(rdb),
	)

	publisher := initNatsPublisher(jsc, logger)

	bus := icommands.NewBus(
		icommands.WithCommand(
			commands.SaveMessageCommand{},
			icommands.NewCommandHandlerMiddleware(
				tracer, commands.NewSaveMessageCommandHandler(idp, stge, publisher, logger)),
		),
	)

	err = replayEventsWithTemporaryConsumer(jsc, bus, tracer, logger)
	if err != nil {
		panic(err)
	}

	msgConsumer, err := initPublishedMessagesNatsConsumer(idp, bus, jsc, tracer, logger)
	go func() {
		if err = msgConsumer.Consume(context.Background(), "rooms.*.messages.*.published"); err != nil {
			panic(err)
		}
	}()

	done := make(chan struct{})
	go func() {
		signal, err := waitForInterrupt()
		if err != nil {
			logger.Error().Err(err).Msg("failed to wait for interrupt")
		}
		logger.Debug().Msgf("received signal: %s", signal)
		done <- struct{}{}
	}()

	<-done
	logger.Info().Msg("messages writer service stopped")
}

func waitForInterrupt() (os.Signal, error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	return <-c, nil
}

func initNatsPublisher(jsc nats.JetStreamContext, logger *zerolog.Logger) *inats.StreamPublisher {
	return inats.NewStreamPublisher(jsc, messagesNatsStream)
}

func initPublishedMessagesNatsConsumer(idp idn.Provider, bus *icommands.Bus, jsc nats.JetStreamContext, tracer tpkg.Tracer, logger *zerolog.Logger) (*inats.PersistentStreamConsumer, error) {
	handlers := pubsub.Handlers{
		events.MessagePublishedEvent{}.SubjectRegex(): events.NewEventHandlerTelemetryMiddleware(tracer, mnats.NewMessagePublishedEventHandler(bus, logger)),
	}

	consumer, err := inats.NewPersistentStreamConsumer(
		jsc,
		handlers,
		logger,
		inats.WithMustHaveHandler(true),
		inats.WithBindStream(messagesNatsStream),
		inats.WithDurable("messages-writer-consumer"),
	)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func replayEventsWithTemporaryConsumer(jsc nats.JetStreamContext, bus *icommands.Bus, tracer tpkg.Tracer, logger *zerolog.Logger) error {
	logger.Debug().Msg("replaying events")

	handlers := pubsub.Handlers{
		events.MessagePublishedEvent{}.SubjectRegex(): events.NewEventHandlerTelemetryMiddleware(
			tracer, mnats.NewMessagePublishedEventHandler(bus, logger)),
	}

	consumer, err := inats.NewPersistentStreamConsumer(
		jsc,
		handlers,
		logger,
		inats.WithMustHaveHandler(true),
		inats.WithCloseOnDone(),
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if consumer.Consume(ctx, "rooms.*.messages.*.published") != nil {
		return err
	}

	logger.Debug().Msg("replaying events done")

	return nil
}

func serveMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
