package commands

import (
	"time"

	icommands "github.com/xfrr/dyschat/internal/commands"
	"github.com/xfrr/dyschat/messages"
)

var _ icommands.Command = (*SaveMessageCommand)(nil)

const (
	SaveMessageCommandType icommands.Type = "save_message"
)

type SaveMessageCommand struct {
	ID          string
	RoomID      string
	UserID      string
	Text        string
	Metadata    messages.Metadata
	PublishedAt time.Time
}

func (SaveMessageCommand) Type() icommands.Type {
	return SaveMessageCommandType
}
