package commands

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/internal/commands"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ commands.CommandHandler = (*BroadcastMessageCommandHandler)(nil)

type BroadcastMessageCommandHandler struct {
	idp       idn.Provider
	stge      agent.RoomStorage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewBroadcastMessageCommandHandler(idp idn.Provider, stge agent.RoomStorage, publisher pubsub.Publisher, logger *zerolog.Logger) *BroadcastMessageCommandHandler {
	return &BroadcastMessageCommandHandler{
		idp:       idp,
		stge:      stge,
		publisher: publisher,
		logger:    logger,
	}
}

func (h *BroadcastMessageCommandHandler) Handle(ctx context.Context, cmd commands.Command) (commands.Reply, error) {
	command, ok := cmd.(*BroadcastMessageCommand)
	if !ok {
		return nil, commands.NewErrInvalidCommandType(cmd.Type())
	}

	room, err := h.stge.Get(ctx, command.RoomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", agent.ErrRoomNotFound, err)
	}

	_, err = room.GetMember(command.UserID)
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(command)
	if err != nil {
		return nil, err
	}

	err = room.Broadcast(ctx, command.UserID, agent.NewEventMessage(
		h.idp.ID(),
		payload,
	))
	if err != nil {
		return nil, err
	}

	h.logger.Debug().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("message broadcasted to room")

	return nil, nil
}
