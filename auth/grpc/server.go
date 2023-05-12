package grpc_api

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/auth"
	proto "github.com/xfrr/dyschat/proto/auth/v1"
)

var _ proto.AuthServiceServer = (*AuthServer)(nil)

type AuthServer struct {
	proto.AuthServiceServer

	authenticator auth.Authenticator
	logger        *zerolog.Logger
}

func NewGrpcServer(auth auth.Authenticator, logger *zerolog.Logger) *AuthServer {
	return &AuthServer{
		authenticator: auth,
		logger:        logger,
	}
}

func (s *AuthServer) Identify(ctx context.Context, req *proto.IdentifyRequest) (*proto.Identity, error) {
	sessionId, err := s.authenticator.Identify(ctx, req.GetToken())
	if err != nil {
		return nil, err
	}

	return &proto.Identity{
		SessionId: sessionId,
	}, nil
}

func (s *AuthServer) Issue(ctx context.Context, req *proto.IssueRequest) (*proto.Key, error) {
	token, err := s.authenticator.Issue(ctx, req.GetSessionId(), req.GetTTL())
	if err != nil {
		return nil, err
	}

	return &proto.Key{
		Token: token,
	}, nil
}
