package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/xfrr/dyschat/auth"
)

type SecretKey string

type JWTParser struct {
	secret SecretKey
}

type ParserOpt func(*JWTParser)

func WithSecretKey(secret string) ParserOpt {
	return func(j *JWTParser) {
		j.secret = SecretKey(secret)
	}
}

func NewJWTParser(opts ...ParserOpt) *JWTParser {
	j := &JWTParser{}
	for _, opt := range opts {
		opt(j)
	}
	return j
}

func (a *JWTParser) Decode(token string) (*auth.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(a.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !t.Valid {
		return nil, errors.New("invalid token provided")
	}

	sessionID, err := t.Claims.GetSubject()
	if err != nil {
		return nil, err
	}

	return &auth.Token{
		SessionID: sessionID,
	}, nil
}

func (a *JWTParser) Encode(sessionID string, ttl int64) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		jwt.RegisteredClaims{
			ID:        sessionID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(ttl) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "dyschat.auth",
			Subject:   "session",
		},
	},
	)

	return t.SignedString([]byte(a.secret))
}
