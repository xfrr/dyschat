package auth

import "context"

type Authenticator interface {
	Identify(ctx context.Context, token string) (string, error)
	Issue(ctx context.Context, sessionID string, ttl int64) (string, error)
}

type authenticator struct {
	parser TokenParser
}

func NewAuthenticator(parser TokenParser) Authenticator {
	return &authenticator{
		parser: parser,
	}
}

func (a *authenticator) Identify(ctx context.Context, token string) (string, error) {
	decodedToken, err := a.decodeToken(token)
	if err != nil {
		return "", err
	}

	return decodedToken.SessionID, nil
}

func (a *authenticator) Issue(ctx context.Context, sessionID string, ttl int64) (string, error) {
	return a.parser.Encode(sessionID, ttl)
}

func (a *authenticator) decodeToken(token string) (*Token, error) {
	return a.parser.Decode(token)
}
