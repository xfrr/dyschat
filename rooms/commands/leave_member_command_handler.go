package commands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/pubsub"
	"github.com/xfrr/dyschat/rooms"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.CommandHandler = (*LeaveMemberCommandHandler)(nil)

type LeaveMemberCommandHandler struct {
	stge      rooms.Storage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewLeaveMemberCommandHandler(stge rooms.Storage, pub pubsub.Publisher, logger *zerolog.Logger) *LeaveMemberCommandHandler {
	return &LeaveMemberCommandHandler{
		stge:      stge,
		publisher: pub,
		logger:    logger,
	}
}

func (h *LeaveMemberCommandHandler) Handle(ctx context.Context, cmd icommands.Command) (reply icommands.Reply, err error) {
	defer func() {
		if err != nil {
			// TODO: publish JoinRoomFailed event
			h.logger.Error().
				Err(err).
				Msg("failed to handle join room command")
		}
	}()

	command, ok := cmd.(*LeaveMemberCommand)
	if !ok {
		return nil, ErrInvalidCommandType
	}

	if err := command.validate(); err != nil {
		return nil, err
	}

	room, err := h.stge.Get(ctx, command.RoomID)
	if err != nil {
		return nil, err
	}

	err = room.RemoveMember(command.UserID)
	if err != nil {
		return nil, err
	}

	if _, err := h.stge.Save(ctx, room); err != nil {
		return nil, err
	}

	if err := h.publisher.Publish(ctx, room.Events()); err != nil {
		return nil, err
	}

	h.logger.Debug().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("member left the room")

	return nil, nil
}
