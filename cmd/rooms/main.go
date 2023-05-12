package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/pkg/env"
	"github.com/xfrr/dyschat/pkg/log"
	"github.com/xfrr/dyschat/proto/auth/v1"
	"github.com/xfrr/dyschat/rooms"
	"github.com/xfrr/dyschat/rooms/commands"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	iauth "github.com/xfrr/dyschat/internal/auth"
	icommands "github.com/xfrr/dyschat/internal/commands"
	inats "github.com/xfrr/dyschat/internal/pubsub/nats"
	rgrpc "github.com/xfrr/dyschat/rooms/api/grpc"
	http "github.com/xfrr/dyschat/rooms/api/http"
	rredis "github.com/xfrr/dyschat/rooms/redis"
)

var (
	// common
	authGRPCAddr = env.Get("DYCHAT_AUTH_GRPC_ADDR", "localhost:50051")
	natsURL      = env.Get("DYCHAT_NATS_URL", "nats://localhost:4222")
	redisAddr    = env.Get("DYCHAT_REDIS_ADDR", "localhost:6379")

	// rooms
	grpcPort = env.Get("DYCHAT_ROOMS_GRPC_PORT", "50052")
	httpPort = env.Get("DYCHAT_ROOMS_HTTP_PORT", "51052")
	logLevel = env.Get("DYCHAT_ROOMS_LOG_LEVEL", "debug")
)

func main() {
	ctx := context.Background()
	logger := log.NewZeroLogger(log.ParseLogLevel(logLevel))

	nc, err := nats.Connect(natsURL)
	if err != nil {
		panic(err)
	}
	defer nc.Drain()

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	streams := inats.NewStreams(
		nats.StreamConfig{
			Name: "room_events",
			Subjects: []string{
				"rooms.>",
			},
		},
	)
	jsc, err := inats.NewJetStreamContext(nc, streams)
	if err != nil {
		panic(err)
	}

	storage := rredis.NewStorage(rdb)

	// room messages publisher
	publisher, err := initNatsRoomMessagesPublisher(jsc, logger)
	if err != nil {
		panic(err)
	}

	idp := idn.NewNanoIDProvider()

	cmdbus := icommands.NewBus(
		icommands.WithCommand(
			commands.CreateRoomCommand{},
			commands.NewCreateRoomCommandHandler(idp, storage, publisher, logger),
		),
		icommands.WithCommand(
			commands.JoinMemberCommand{},
			commands.NewJoinMemberCommandHandler(storage, publisher, logger),
		),
		icommands.WithCommand(
			commands.LeaveMemberCommand{},
			commands.NewLeaveMemberCommandHandler(storage, publisher, logger),
		),
	)

	view := rredis.NewView(rdb)
	querier := rooms.NewQuerier(view)

	authGrpcClient, err := newAuthGRPCConn()
	if err != nil {
		panic(err)
	}

	authInterceptor := iauth.NewAuthGRPCInterceptor(authGrpcClient)

	grpcServer := rgrpc.NewRoomsServer(cmdbus, querier, authInterceptor, logger)

	httpAddr := ":" + httpPort
	grpcAddr := ":" + grpcPort

	go func() {
		err := grpcServer.Run(ctx, rgrpc.Addr(grpcAddr))
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		httpServer := http.NewServer(httpAddr, grpcAddr, logger)
		err := httpServer.Serve(ctx)
		if err != nil {
			panic(err)
		}
	}()

	done := make(chan struct{})
	go func() {
		waitForInterrupt()
		done <- struct{}{}
	}()

	<-done
	logger.Info().Msg("rooms service stopped")
}

func waitForInterrupt() (os.Signal, error) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	return <-c, nil
}

func initNatsRoomMessagesPublisher(jsc nats.JetStreamContext, logger *zerolog.Logger) (*inats.StreamPublisher, error) {
	return inats.NewStreamPublisher(jsc, ""), nil
}

func newAuthGRPCConn() (auth.AuthServiceClient, error) {
	conn, err := grpc.Dial(authGRPCAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return auth.NewAuthServiceClient(conn), nil
}
