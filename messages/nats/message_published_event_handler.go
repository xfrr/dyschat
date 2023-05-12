package nats

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/events"
	"github.com/xfrr/dyschat/messages/commands"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

type MessagePublishedEventHandler struct {
	cmdbus *icommands.Bus
	logger *zerolog.Logger
}

func NewMessagePublishedEventHandler(cmdbus *icommands.Bus, logger *zerolog.Logger) *MessagePublishedEventHandler {
	return &MessagePublishedEventHandler{
		cmdbus: cmdbus,
		logger: logger,
	}
}

func (h *MessagePublishedEventHandler) Handle(ctx context.Context, _ string, data []byte) error {
	ev := events.MessagePublishedEvent{}
	if err := ev.UnmarshalJSON(data); err != nil {
		return err
	}

	h.logger.Info().
		Str("message_id", ev.MessageID).
		Str("room_id", ev.RoomID).
		Str("user_id", ev.UserID).
		Msg("received message published event")

	cmd := &commands.SaveMessageCommand{
		ID:          ev.MessageID,
		RoomID:      ev.RoomID,
		UserID:      ev.UserID,
		Text:        ev.Text,
		Metadata:    ev.Metadata,
		PublishedAt: time.Unix(ev.PublishedAt, 0),
	}

	if _, err := h.cmdbus.Dispatch(ctx, cmd); err != nil {
		return err
	}

	return nil
}
