package commands

import (
	"errors"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*JoinMemberCommand)(nil)

const (
	JoinMemberCommandType icommands.Type = "join_room_member"
)

var (
	ErrInvalidSecretKey = errors.New("invalid secret key")
)

type JoinMemberCommand struct {
	RoomID    string `json:"-"`
	UserID    string `json:"user_id"`
	SecretKey string `json:"secret_key"`
}

func (JoinMemberCommand) Type() icommands.Type {
	return JoinMemberCommandType
}
