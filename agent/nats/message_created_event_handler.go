package nats

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent/commands"
	"github.com/xfrr/dyschat/internal/events"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

type MessageSavedEventHandler struct {
	logger *zerolog.Logger

	cmdbus *icommands.Bus
}

func NewMessageSavedEventHandler(cmdbus *icommands.Bus, logger *zerolog.Logger) *MessageSavedEventHandler {
	return &MessageSavedEventHandler{
		logger: logger,
		cmdbus: cmdbus,
	}
}

func (h *MessageSavedEventHandler) Handle(ctx context.Context, _ string, data []byte) error {
	ev := events.MessageCreated{}
	if err := ev.UnmarshalJSON(data); err != nil {
		return err
	}

	h.logger.Debug().
		Any("event", ev).
		Msg("message created event received")

	_, err := h.cmdbus.Dispatch(ctx, &commands.BroadcastMessageCommand{
		MessageID:   ev.MessageID,
		RoomID:      ev.RoomID,
		UserID:      ev.UserID,
		Text:        ev.Text,
		Metadata:    ev.Metadata,
		CreatedAt:   ev.CreatedAt,
		PublishedAt: ev.PublishedAt,
	})
	if err != nil {
		return err
	}

	return nil
}
