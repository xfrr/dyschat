package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*AuthenticateCommand)(nil)

const (
	AuthenticateCommandType icommands.Type = "authenticate"
)

type AuthenticateCommand struct {
	RoomID string
	UserID string
	Token  string
}

func (AuthenticateCommand) Type() icommands.Type {
	return AuthenticateCommandType
}

func (cmd *AuthenticateCommand) validate() error {
	if cmd.RoomID == "" {
		return ErrInvalidRoomID
	}
	if cmd.UserID == "" {
		return ErrInvalidUserID
	}
	if cmd.Token == "" {
		return ErrInvalidSecretKey
	}
	return nil
}
