package grpc_api

import (
	"context"
	"net"

	"github.com/xfrr/dyschat/proto/auth/v1"
	"google.golang.org/grpc"
)

type Addr string

func (s *AuthServer) Run(ctx context.Context, addr Addr) error {
	listener, err := net.Listen("tcp", string(addr))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	auth.RegisterAuthServiceServer(server, s)

	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	s.logger.Debug().
		Str("addr", string(addr)).
		Msg("starting grpc server")

	return server.Serve(listener)
}
