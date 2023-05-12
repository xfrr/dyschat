package commands

import (
	"errors"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

var (
	ErrInvalidCommandType = icommands.NewErrInvalidCommandType("publish_message")

	ErrInvalidUserID      = errors.New("invalid user id")
	ErrorInvalidSecretKey = errors.New("invalid secret key")

	ErrInvalidRoomID    = errors.New("invalid room id")
	ErrInvalidRoomName  = errors.New("invalid room name")
	ErrInvalidMemberIDs = errors.New("invalid member ids")
	ErrInvalidCreatedAt = errors.New("invalid created at")
)
