package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/agent/commands"
	"github.com/xfrr/dyschat/agent/inmemory"
	"github.com/xfrr/dyschat/agent/ws"
	"github.com/xfrr/dyschat/internal/events"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
	"github.com/xfrr/dyschat/pkg/env"
	"github.com/xfrr/dyschat/pkg/log"
	"github.com/xfrr/dyschat/pkg/telemetry"
	"github.com/xfrr/dyschat/pkg/telemetry/prometheus"

	hnats "github.com/xfrr/dyschat/agent/nats"
	icommands "github.com/xfrr/dyschat/internal/commands"
	inats "github.com/xfrr/dyschat/internal/pubsub/nats"
	xprom "go.opentelemetry.io/otel/exporters/prometheus"
)

var (
	svcName = "ws-agent"

	natsURL = env.Get("DYSCHAT_NATS_URL", "nats://0.0.0.0:4222")

	wsPort   = env.Get("DYSCHAT_WS_AGENT_PORT", "8083")
	logLevel = env.Get("DYSCHAT_WS_AGENT_LOG_LEVEL", "debug")
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
			Name:      "messages",
			Subjects:  []string{"rooms.*.messages.>"},
			Retention: nats.InterestPolicy,
		},
		nats.StreamConfig{
			Name: "room_events",
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

	storage := inmemory.NewRoomStorage()

	err = replayNatsRoomEventsStreamConsumer(idp, storage, jsc, logger)
	if err != nil {
		panic(err)
	}

	// messages consumer
	consumer, err := initNatsMessagesStreamConsumer(idp, storage, jsc, logger)
	go func() {
		if err = consumer.Consume(context.Background(), "rooms.*.messages.>"); err != nil {
			panic(err)
		}
	}()

	// room events consumer
	reconsumer, err := initNatsRoomEventsStreamConsumer(idp, storage, jsc, logger)
	go func() {
		if err = reconsumer.Consume(context.Background(), "rooms.>"); err != nil {
			panic(err)
		}
	}()

	publisher, err := initNatsRoomMessagesPublisher(jsc, logger)
	if err != nil {
		panic(err)
	}

	exporter, err := xprom.New()
	if err != nil {
		panic(err)
	}

	meter := prometheus.NewMeterProvider(svcName, exporter)
	go serveMetrics()

	wsServer := initWebsocketServer(idp, storage, publisher, meter, logger)
	var done = make(chan struct{})

	go func() {
		if err = wsServer.Start(context.Background()); err != nil {
			panic(err)
		}
	}()

	go func() {
		signal, err := waitForInterrupt()
		if err != nil {
			logger.Error().Err(err).Msg("failed to wait for interrupt")
		}
		logger.Debug().Msgf("received signal: %s", signal)
		done <- struct{}{}
	}()

	<-done
	logger.Info().Msg("websocket agent stopped")
}

func waitForInterrupt() (os.Signal, error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	return <-c, nil
}

func initNatsRoomEventsStreamConsumer(idp idn.Provider, stge agent.RoomStorage, jsc nats.JetStreamContext, logger *zerolog.Logger) (*inats.PersistentStreamConsumer, error) {
	bus := icommands.NewBus(
		icommands.WithCommand(
			commands.CreateRoomCommand{},
			commands.NewCreateRoomCommandHandler(stge, logger),
		),
		icommands.WithCommand(
			commands.JoinRoomCommand{},
			commands.NewJoinRoomCommandHandler(idp, stge, nil, logger),
		),
	)

	return inats.NewPersistentStreamConsumer(jsc, pubsub.Handlers{
		events.RoomCreatedEvent{}.SubjectRegex(): hnats.NewRoomCreatedEventHandler(bus, logger),
		events.RoomMemberJoined{}.SubjectRegex(): hnats.NewMemberJoinedEventHandler(bus, logger),
	},
		logger,
		inats.WithMustHaveHandler(false),
		inats.WithBindStream("room_events"),
		inats.WithDurable("room_events_ws_agent_consumer"),
	)
}

func initNatsMessagesStreamConsumer(idp idn.Provider, stge agent.RoomStorage, jsc nats.JetStreamContext, logger *zerolog.Logger) (*inats.PersistentStreamConsumer, error) {
	bus := icommands.NewBus(
		icommands.WithCommand(
			commands.BroadcastMessageCommand{},
			commands.NewBroadcastMessageCommandHandler(idp, stge, nil, logger),
		),
	)

	return inats.NewPersistentStreamConsumer(jsc, pubsub.Handlers{
		events.MessageCreated{}.SubjectRegex(): hnats.NewMessageSavedEventHandler(bus, logger),
	},
		logger,
		inats.WithMustHaveHandler(false),
		inats.WithBindStream("messages"),
		inats.WithDurable("messages_ws_agent_consumer"),
	)
}

func replayNatsRoomEventsStreamConsumer(idp idn.Provider, stge agent.RoomStorage, jsc nats.JetStreamContext, logger *zerolog.Logger) error {
	bus := icommands.NewBus(
		icommands.WithCommand(
			commands.CreateRoomCommand{},
			commands.NewCreateRoomCommandHandler(stge, logger),
		),
		icommands.WithCommand(
			commands.JoinRoomCommand{},
			commands.NewJoinRoomCommandHandler(idp, stge, nil, logger),
		),
		icommands.WithCommand(
			commands.BroadcastMessageCommand{},
			commands.NewBroadcastMessageCommandHandler(idp, stge, nil, logger),
		),
	)

	consumer, err := inats.NewPersistentStreamConsumer(jsc, pubsub.Handlers{
		events.RoomCreatedEvent{}.SubjectRegex(): hnats.NewRoomCreatedEventHandler(bus, logger),
		events.RoomMemberJoined{}.SubjectRegex(): hnats.NewMemberJoinedEventHandler(bus, logger),
		events.MessageCreated{}.SubjectRegex():   hnats.NewMessageSavedEventHandler(bus, logger),
	},
		logger,
		inats.WithMustHaveHandler(true),
		inats.WithCloseOnDone(),
	)
	if err != nil {
		return err
	}

	logger.Debug().Msg("replaying events")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err = consumer.Consume(ctx, "rooms.*.events.>"); err != nil {
		return err
	}

	if err = consumer.Consume(ctx, "rooms.*.members.>"); err != nil {
		return err
	}

	if err = consumer.Consume(ctx, "rooms.*.messages.*.created"); err != nil {
		return err
	}

	logger.Debug().Msg("events replayed")
	return nil
}

func initNatsRoomMessagesPublisher(jsc nats.JetStreamContext, logger *zerolog.Logger) (*inats.StreamPublisher, error) {
	return inats.NewStreamPublisher(jsc, "messages"), nil
}

func initWebsocketServer(idp idn.Provider, stge agent.RoomStorage, pub pubsub.Publisher, meter telemetry.Meter, logger *zerolog.Logger) *ws.Server {
	bus := icommands.NewBus(
		icommands.WithCommand(
			commands.AuthenticateCommand{},
			commands.NewAuthenticateCommandHandler(stge, logger),
		),
		icommands.WithCommand(
			commands.ConnectCommand{},
			commands.NewConnectCommandHandler(idp, stge, pub, logger),
		),
		icommands.WithCommand(
			commands.PublishMessageCommand{},
			commands.NewPublishMessageCommandHandler(idp, stge, pub, logger),
		),
	)

	return ws.NewServer(
		bus,
		logger,
		ws.WithAddr(":"+wsPort),
		ws.WithMeter(meter),
	)
}

func serveMetrics() {
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":7777", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
