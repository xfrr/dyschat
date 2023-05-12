package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*JoinRoomCommand)(nil)

const (
	JoinRoomCommandType icommands.Type = "join_room"
)

type JoinRoomCommand struct {
	UserID string
	RoomID string
}

func (JoinRoomCommand) Type() icommands.Type {
	return JoinRoomCommandType
}

func (cmd *JoinRoomCommand) validate() error {
	if cmd.RoomID == "" {
		return ErrInvalidRoomID
	}
	if cmd.UserID == "" {
		return ErrInvalidUserID
	}

	return nil
}
