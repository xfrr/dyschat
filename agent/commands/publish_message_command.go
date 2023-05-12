package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*PublishMessageCommand)(nil)

const (
	PublishMessageCommandType icommands.Type = "message"
)

type PublishMessageCommand struct {
	RoomID string
	UserID string
	Text   string
}

func NewPublishMessageCommand(userID, roomID string, text string) (*PublishMessageCommand, error) {
	cmd := &PublishMessageCommand{
		UserID: userID,
		RoomID: roomID,
		Text:   text,
	}

	return cmd, cmd.validate()
}

func (cmd *PublishMessageCommand) validate() error {
	if cmd.UserID == "" {
		return ErrInvalidUserID
	}

	if cmd.RoomID == "" {
		return ErrInvalidRoomID
	}

	if cmd.Text == "" {
		return ErrInvalidMessageText
	}

	return nil
}

func (PublishMessageCommand) Type() icommands.Type {
	return PublishMessageCommandType
}
