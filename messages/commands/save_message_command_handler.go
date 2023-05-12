package commands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
	"github.com/xfrr/dyschat/messages"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ icommands.CommandHandler = (*SaveMessageCommandHandler)(nil)

type SaveMessageCommandHandler struct {
	idp       idn.Provider
	stge      messages.Storage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewSaveMessageCommandHandler(idp idn.Provider, stge messages.Storage, pub pubsub.Publisher, logger *zerolog.Logger) *SaveMessageCommandHandler {
	return &SaveMessageCommandHandler{
		idp:       idp,
		stge:      stge,
		publisher: pub,
		logger:    logger,
	}
}

func (h *SaveMessageCommandHandler) Handle(ctx context.Context, cmd icommands.Command) (reply icommands.Reply, err error) {
	defer func() {
		if err != nil {
			h.logger.Error().
				Err(err).
				Msg("failed saving message")
		}
	}()

	command, ok := cmd.(*SaveMessageCommand)
	if !ok {
		return nil, ErrInvalidCommandType
	}

	message, err := messages.NewMessage(
		command.ID,
		command.RoomID,
		command.UserID,
		command.Text,
		messages.WithMetadata(command.Metadata),
		messages.WithPublishedAt(command.PublishedAt))
	if err != nil {
		return nil, err
	}

	h.logger.Debug().
		Str("id", message.ID()).
		Str("room_id", message.RoomID()).
		Str("user_id", message.UserID()).
		Str("text", message.Text()).
		Msg("saving message")

	if err := h.stge.Save(ctx, message); err != nil {
		return nil, err
	}

	if err := h.publisher.Publish(ctx, message.Events()); err != nil {
		return nil, err
	}

	h.logger.Debug().
		Str("id", message.ID()).
		Str("room_id", message.RoomID()).
		Str("user_id", message.UserID()).
		Str("text", message.Text()).
		Msg("message saved")

	return nil, nil
}
