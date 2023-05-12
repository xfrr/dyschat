package commands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
	"github.com/xfrr/dyschat/rooms"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.CommandHandler = (*CreateRoomCommandHandler)(nil)

type CreateRoomReply struct {
	RoomID    string `json:"room_id"`
	SecretKey string `json:"secret_key"`
}

type CreateRoomCommandHandler struct {
	idp       idn.Provider
	stge      rooms.Storage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewCreateRoomCommandHandler(idp idn.Provider, stge rooms.Storage, pub pubsub.Publisher, logger *zerolog.Logger) *CreateRoomCommandHandler {
	return &CreateRoomCommandHandler{
		idp:       idp,
		stge:      stge,
		publisher: pub,
		logger:    logger,
	}
}

func (h *CreateRoomCommandHandler) Handle(ctx context.Context, cmd icommands.Command) (reply icommands.Reply, err error) {
	defer func() {
		if err != nil {
			h.logger.Error().
				Err(err).
				Msg("failed creating room")
		}
	}()

	command, ok := cmd.(*CreateRoomCommand)
	if !ok {
		return nil, ErrInvalidCommandType
	}

	_, err = h.stge.Get(ctx, command.ID)
	if err != nil && err != rooms.ErrRoomNotFound {
		return nil, err
	}

	room, err := rooms.NewRoom(command.ID, command.Name, command.OwnerID, h.idp.ID())
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
		Str("id", room.ID()).
		Str("name", room.Name()).
		Msg("room created")

	return &CreateRoomReply{
		RoomID:    command.ID,
		SecretKey: room.SecretKey(),
	}, nil
}
