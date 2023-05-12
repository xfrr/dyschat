package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*CreateRoomCommand)(nil)

const (
	CreateRoomCommandType icommands.Type = "create_room"
)

type CreateRoomCommand struct {
	ID        string `json:"-"`
	Name      string `json:"name"`
	SecretKey string `json:"secret_key"`
	OwnerID   string `json:"owner_id"`
}

func (CreateRoomCommand) Type() icommands.Type {
	return CreateRoomCommandType
}
