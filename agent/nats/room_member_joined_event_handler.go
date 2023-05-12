package nats

import (
	"context"

	"github.com/rs/zerolog"

	"github.com/xfrr/dyschat/agent/commands"
	icommands "github.com/xfrr/dyschat/internal/commands"
	"github.com/xfrr/dyschat/internal/events"
)

type MemberJoinedEventHandler struct {
	cmdbus *icommands.Bus
	logger *zerolog.Logger
}

func NewMemberJoinedEventHandler(bus *icommands.Bus, logger *zerolog.Logger) *MemberJoinedEventHandler {
	return &MemberJoinedEventHandler{
		cmdbus: bus,
		logger: logger,
	}
}

func (h *MemberJoinedEventHandler) Handle(ctx context.Context, _ string, data []byte) error {
	h.logger.Debug().
		RawJSON("data", data).
		Msg("member joined room event received")

	ev := events.RoomMemberJoined{}
	if err := ev.UnmarshalJSON(data); err != nil {
		return err
	}

	cmd := &commands.JoinRoomCommand{
		RoomID: ev.RoomID,
		UserID: ev.UserID,
	}

	_, err := h.cmdbus.Dispatch(ctx, cmd)
	if err != nil {
		return err
	}

	return nil
}
