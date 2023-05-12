package nats

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent/commands"
	"github.com/xfrr/dyschat/internal/events"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

type RoomCreatedEventHandler struct {
	logger *zerolog.Logger

	cmdbus *icommands.Bus
}

func NewRoomCreatedEventHandler(cmdbus *icommands.Bus, logger *zerolog.Logger) *RoomCreatedEventHandler {
	return &RoomCreatedEventHandler{
		logger: logger,
		cmdbus: cmdbus,
	}
}

func (h *RoomCreatedEventHandler) Handle(ctx context.Context, _ string, data []byte) error {
	ev := events.RoomCreatedEvent{}
	if err := ev.UnmarshalJSON(data); err != nil {
		return err
	}

	h.logger.Debug().
		Any("event", ev).
		Msg("room created event received")

	_, err := h.cmdbus.Dispatch(ctx, &commands.CreateRoomCommand{
		RoomID:    ev.RoomID,
		SecretKey: ev.SecretKey,
		CreatedAt: ev.CreatedAt,
	})
	if err != nil {
		return err
	}

	return nil
}
