package commands

import (
	"context"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/internal/commands"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ commands.CommandHandler = (*JoinRoomCommandHandler)(nil)

type JoinRoomCommandHandler struct {
	idp       idn.Provider
	stge      agent.RoomStorage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewJoinRoomCommandHandler(idp idn.Provider, stge agent.RoomStorage, pub pubsub.Publisher, logger *zerolog.Logger) *JoinRoomCommandHandler {
	return &JoinRoomCommandHandler{
		idp:       idp,
		stge:      stge,
		publisher: pub,
		logger:    logger,
	}
}

func (h *JoinRoomCommandHandler) Handle(ctx context.Context, cmd commands.Command) (commands.Reply, error) {
	command, ok := cmd.(*JoinRoomCommand)
	if !ok {
		return nil, commands.NewErrInvalidCommandType(cmd.Type())
	}
	if err := command.validate(); err != nil {
		return nil, err
	}

	room, err := h.stge.Get(ctx, command.RoomID)
	if err != nil {
		return nil, err
	}

	_, err = room.AddMember(command.UserID)
	if err != nil {
		return nil, err
	}

	err = h.stge.Save(ctx, room)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	msg := agent.NewEventMessage(h.idp.ID(), payload)
	return nil, room.Broadcast(ctx, "system", msg)
}
