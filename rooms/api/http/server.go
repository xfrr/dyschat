package http

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/proto/rooms/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

type Server struct {
	name     string
	httpAddr string
	grpcAddr string
	log      *zerolog.Logger
}

func NewServer(httpAddr, grpcAddr string, log *zerolog.Logger) Server {
	return Server{name: "rooms-http-server", httpAddr: httpAddr, grpcAddr: grpcAddr, log: log}
}

func (s *Server) Serve(ctx context.Context) error {
	// TODO: Implements TLS

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))

	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	}

	// Register Prometheus metrics handler.
	mux.HandlePath("GET", "/metrics", s.promHandler)

	err := rooms.RegisterRoomsServiceHandlerFromEndpoint(ctx, mux, s.grpcAddr, opts)
	if err != nil {
		return err
	}

	s.log.Info().
		Str("addr", s.httpAddr).
		Msg("starting http server")

	return http.ListenAndServe(s.httpAddr, mux)
}

func (s *Server) promHandler(w http.ResponseWriter, r *http.Request, params map[string]string) {
	h := promhttp.Handler()
	h.ServeHTTP(w, r)
}
