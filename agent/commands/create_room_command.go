package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*CreateRoomCommand)(nil)

const (
	CreateRoomCommandType icommands.Type = "create_room"
)

type CreateRoomCommand struct {
	RoomID    string
	SecretKey string
	CreatedAt int64
}

func (cmd *CreateRoomCommand) validate() error {
	if cmd.RoomID == "" {
		return ErrInvalidRoomID
	}

	if cmd.SecretKey == "" {
		return ErrInvalidSecretKey
	}

	if cmd.CreatedAt == 0 {
		return ErrInvalidCreatedAt
	}

	return nil
}

func (CreateRoomCommand) Type() icommands.Type {
	return CreateRoomCommandType
}
