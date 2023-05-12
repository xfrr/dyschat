package grpc

import (
	"context"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	ggrpc "google.golang.org/grpc"
)

type Server struct {
	cfg        *ServerConfig
	grpcServer *ggrpc.Server
}

type ServerConfig struct {
	Addr string

	iunary  []ggrpc.UnaryServerInterceptor
	istream []ggrpc.StreamServerInterceptor
}

type ServerOption func(*ServerConfig)

func WithUnaryServerInterceptor(i ...ggrpc.UnaryServerInterceptor) ServerOption {
	return func(c *ServerConfig) {
		c.iunary = append(c.iunary, i...)
	}
}

func WithStreamServerInterceptor(i ...ggrpc.StreamServerInterceptor) ServerOption {
	return func(c *ServerConfig) {
		c.istream = append(c.istream, i...)
	}
}

func NewServer(ctx context.Context, addr string, opts ...ServerOption) (*Server, error) {
	cfg := &ServerConfig{
		Addr: addr,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &Server{
		cfg: cfg,
		grpcServer: ggrpc.NewServer(
			ggrpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(cfg.iunary...)),
			ggrpc.StreamInterceptor(grpc_middleware.ChainStreamServer(cfg.istream...)),
		),
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return nil
}

func (s *Server) RegisterService(sd *ggrpc.ServiceDesc, ss interface{}) {
	s.grpcServer.RegisterService(sd, ss)
}
