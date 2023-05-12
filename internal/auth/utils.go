package auth

import (
	"context"
	"errors"
)

func UserIDFromContext(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDMetadataKey).(string)
	if !ok {
		return "", errors.New("user id not found in context")
	}

	return userID, nil
}
