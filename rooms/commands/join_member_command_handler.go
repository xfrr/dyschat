package commands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/pubsub"
	"github.com/xfrr/dyschat/rooms"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.CommandHandler = (*JoinMemberCommandHandler)(nil)

type JoinMemberCommandHandler struct {
	stge      rooms.Storage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewJoinMemberCommandHandler(stge rooms.Storage, pub pubsub.Publisher, logger *zerolog.Logger) *JoinMemberCommandHandler {
	return &JoinMemberCommandHandler{
		stge:      stge,
		publisher: pub,
		logger:    logger,
	}
}

func (h *JoinMemberCommandHandler) Handle(ctx context.Context, cmd icommands.Command) (reply icommands.Reply, err error) {
	defer func() {
		if err != nil {
			// TODO: publish JoinRoomFailed event
			h.logger.Error().
				Err(err).
				Msg("failed to handle join room command")
		}
	}()

	command, ok := cmd.(*JoinMemberCommand)
	if !ok {
		return nil, ErrInvalidCommandType
	}

	room, err := h.stge.Get(ctx, command.RoomID)
	if err != nil {
		return nil, err
	}

	if !room.IsSecretKeyValid(command.SecretKey) {
		return nil, rooms.ErrUnauthorized
	}

	_, err = room.AddMember(command.UserID)
	if err != nil {
		return nil, err
	}

	if err := h.publisher.Publish(ctx, room.Events()); err != nil {
		return nil, err
	}

	if _, err := h.stge.Save(ctx, room); err != nil {
		return nil, err
	}

	h.logger.Debug().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("member joined the room")

	return nil, nil
}
