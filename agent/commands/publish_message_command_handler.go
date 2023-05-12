package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/internal/commands"
	"github.com/xfrr/dyschat/internal/events"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ commands.CommandHandler = (*PublishMessageCommandHandler)(nil)

type PublishMessageCommandHandler struct {
	idp       idn.Provider
	stge      agent.RoomStorage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewPublishMessageCommandHandler(idp idn.Provider, stge agent.RoomStorage, publisher pubsub.Publisher, logger *zerolog.Logger) *PublishMessageCommandHandler {
	return &PublishMessageCommandHandler{
		idp:       idp,
		stge:      stge,
		publisher: publisher,
		logger:    logger,
	}
}

func (h *PublishMessageCommandHandler) Handle(ctx context.Context, cmd commands.Command) (commands.Reply, error) {
	command, ok := cmd.(*PublishMessageCommand)
	if !ok {
		return nil, commands.NewErrInvalidCommandType(cmd.Type())
	}

	if err := command.validate(); err != nil {
		return nil, err
	}

	room, err := h.stge.Get(ctx, command.RoomID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", agent.ErrRoomNotFound, err)
	}

	member, err := room.GetMember(command.UserID)
	if err != nil {
		return nil, err
	}

	if !member.IsAuthenticated() {
		return nil, agent.ErrUnauthenticated
	}

	event := &events.MessagePublishedEvent{
		MessageID:   h.idp.ID(),
		Text:        command.Text,
		RoomID:      room.ID(),
		UserID:      member.ID(),
		PublishedAt: time.Now().Unix(),
	}

	if err := h.publisher.Publish(ctx, []pubsub.Event{event}); err != nil {
		return nil, err
	}

	h.logger.Debug().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("message published to room")

	return nil, nil
}
