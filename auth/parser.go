package auth

type TokenParser interface {
	Decode(token string) (*Token, error)
	Encode(sessionID string, ttl int64) (string, error)
}
