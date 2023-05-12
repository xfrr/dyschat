package commands

import (
	"errors"
)

var (
	ErrInvalidUserID    = errors.New("invalid user id")
	ErrInvalidSecretKey = errors.New("invalid secret key")

	ErrInvalidRoomID    = errors.New("invalid room id")
	ErrInvalidMemberIDs = errors.New("invalid member ids")
	ErrInvalidCreatedAt = errors.New("invalid created at")

	ErrInvalidMessageID      = errors.New("invalid message id")
	ErrInvalidMessage        = errors.New("invalid message received")
	ErrInvalidMessagePayload = errors.New("invalid message payload")
	ErrInvalidMessageText    = errors.New("invalid message text")
)
