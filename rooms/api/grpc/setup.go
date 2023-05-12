package grpc

import (
	"context"
	"net"

	"github.com/xfrr/dyschat/proto/rooms/v1"
	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

type Addr string

func (s *RoomsServer) Run(ctx context.Context, addr Addr) error {
	listener, err := net.Listen("tcp", string(addr))
	if err != nil {
		return err
	}

	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			s.authInterceptor.Unary(),
		),
	)

	rooms.RegisterRoomsServiceServer(server, s)

	go func() {
		<-ctx.Done()
		server.Stop()
	}()

	s.logger.Debug().
		Str("addr", string(addr)).
		Msg("starting grpc server")

	return server.Serve(listener)
}
