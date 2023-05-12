package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*LeaveMemberCommand)(nil)

const (
	LeaveMemberCommandType icommands.Type = "leave_room_member"
)

type LeaveMemberCommand struct {
	RoomID string
	UserID string
}

func (LeaveMemberCommand) Type() icommands.Type {
	return LeaveMemberCommandType
}

func (cmd *LeaveMemberCommand) validate() error {
	if cmd.RoomID == "" {
		return ErrInvalidRoomID
	}

	if cmd.UserID == "" {
		return ErrInvalidUserID
	}

	return nil
}
