package messages

import "errors"

var (
	ErrInvalidMessageID   = errors.New("invalid message id")
	ErrInvalidRoomID      = errors.New("invalid room id")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrInvalidText        = errors.New("invalid text")
	ErrInvalidMetadata    = errors.New("invalid metadata")
	ErrInvalidCreatedAt   = errors.New("invalid created at")
	ErrInvalidPublishedAt = errors.New("invalid published at")

	ErrInvalidMessage = errors.New("invalid message")

	ErrMessageNotFound = errors.New("message not found")
)
