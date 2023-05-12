package commands

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/internal/commands"
)

var _ commands.CommandHandler = (*CreateRoomCommandHandler)(nil)

type CreateRoomCommandHandler struct {
	stge agent.RoomStorage

	logger *zerolog.Logger
}

func NewCreateRoomCommandHandler(stge agent.RoomStorage, logger *zerolog.Logger) *CreateRoomCommandHandler {
	return &CreateRoomCommandHandler{
		stge:   stge,
		logger: logger,
	}
}

func (h *CreateRoomCommandHandler) Handle(ctx context.Context, cmd commands.Command) (commands.Reply, error) {
	command, ok := cmd.(*CreateRoomCommand)
	if !ok {
		return nil, commands.NewErrInvalidCommandType(cmd.Type())
	}

	if err := command.validate(); err != nil {
		return nil, err
	}

	_, err := h.stge.Get(ctx, command.RoomID)
	if err == nil {
		return nil, err
	}

	createdAt := time.Unix(command.CreatedAt, 0)
	room, err := agent.NewRoom(command.RoomID, command.SecretKey, createdAt)
	if err != nil {
		return nil, err
	}

	err = h.stge.Save(ctx, room)
	if err != nil {
		return nil, err
	}

	h.logger.Debug().
		Str("room_id", command.RoomID).
		Int64("created_at", command.CreatedAt).
		Msg("created room")
	return nil, nil
}
