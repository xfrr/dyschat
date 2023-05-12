package commands

import (
	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.Command = (*BroadcastMessageCommand)(nil)

const (
	BroadcastMessageCommandType icommands.Type = "broadcast_message"
)

type BroadcastMessageCommand struct {
	MessageID   string                 `json:"message_id"`
	RoomID      string                 `json:"room_id"`
	UserID      string                 `json:"user_id"`
	Text        string                 `json:"text"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   int64                  `json:"created_at"`
	PublishedAt int64                  `json:"published_at"`
}

func (BroadcastMessageCommand) Type() icommands.Type {
	return BroadcastMessageCommandType
}
