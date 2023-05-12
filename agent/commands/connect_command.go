package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*ConnectCommand)(nil)

const (
	ConnectCommandType icommands.Type = "connect"
)

type ConnectCommand struct {
	RoomID    string `json:"-"`
	UserID    string `json:"user_id"`
	SecretKey string `json:"secret_key"`

	Send chan []byte `json:"-"`
}

func (ConnectCommand) Type() icommands.Type {
	return ConnectCommandType
}

func (cmd *ConnectCommand) validate() error {
	if cmd.RoomID == "" {
		return ErrInvalidRoomID
	}
	if cmd.UserID == "" {
		return ErrInvalidUserID
	}

	return nil
}
