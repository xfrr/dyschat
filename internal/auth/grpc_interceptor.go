package auth

import (
	"context"

	"github.com/xfrr/dyschat/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type SessionKey string

const (
	UserIDMetadataKey SessionKey = "user_id"
)

type AuthGRPCInterceptor struct {
	authClient auth.AuthServiceClient
}

func NewAuthGRPCInterceptor(authClient auth.AuthServiceClient) *AuthGRPCInterceptor {
	return &AuthGRPCInterceptor{authClient: authClient}
}

func (i *AuthGRPCInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// get token from metadata
		token, err := i.getTokenFromMetadata(ctx)
		if err != nil {
			return nil, err
		}

		// validate token
		res, err := i.authClient.Identify(ctx, &auth.IdentifyRequest{Token: token})
		if err != nil {
			return nil, err
		}

		// add user id to context
		ctx = context.WithValue(ctx, UserIDMetadataKey, res.GetSessionId())

		return handler(ctx, req)
		// return nil, nil
	}
}

func (i *AuthGRPCInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	}
}

func (i *AuthGRPCInterceptor) getTokenFromMetadata(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	if len(md["authorization"]) != 1 {
		return "", status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return md["authorization"][0], nil
}
