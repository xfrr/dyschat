package main

import (
	"context"

	"github.com/xfrr/dyschat/auth"
	grpc_api "github.com/xfrr/dyschat/auth/grpc"
	"github.com/xfrr/dyschat/auth/jwt"
	"github.com/xfrr/dyschat/pkg/env"
	"github.com/xfrr/dyschat/pkg/log"
)

var (
	grpcPort = env.Get("DYSCHAT_AUTH_GRPC_PORT", "50051")
	secret   = env.Get("DYSCHAT_AUTH_JWT_SECRET", "secret")
	logLevel = env.Get("DYSCHAT_AUTH_LOG_LEVEL", "debug")
)

func main() {
	ctx := context.Background()
	logger := log.NewZeroLogger(log.ParseLogLevel(logLevel))
	authenticator := auth.NewAuthenticator(jwt.NewJWTParser(jwt.WithSecretKey(secret)))
	server := grpc_api.NewGrpcServer(authenticator, logger).Run(ctx, grpc_api.Addr(":"+grpcPort))
	if err := server; err != nil {
		logger.Fatal().Err(err).Msg("failed to start grpc server")
		return
	}
}
