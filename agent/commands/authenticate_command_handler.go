package commands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/internal/commands"
)

var _ commands.CommandHandler = (*AuthenticateCommandHandler)(nil)

type AuthenticateCommandHandler struct {
	stge agent.RoomStorage

	logger *zerolog.Logger
}

func NewAuthenticateCommandHandler(stge agent.RoomStorage, logger *zerolog.Logger) *AuthenticateCommandHandler {
	return &AuthenticateCommandHandler{
		stge:   stge,
		logger: logger,
	}
}

func (h *AuthenticateCommandHandler) Handle(ctx context.Context, cmd commands.Command) (_ commands.Reply, err error) {
	command, ok := cmd.(*AuthenticateCommand)
	if !ok {
		return nil, commands.NewErrInvalidCommandType(cmd.Type())
	}

	h.logger.Debug().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("authenticating user")

	if err = command.validate(); err != nil {
		return nil, err
	}

	room, err := h.stge.Get(ctx, command.RoomID)
	if err != nil {
		return nil, err
	}

	member, err := room.GetMember(command.UserID)
	if err != nil {
		return nil, err
	}

	if member.IsAuthenticated() {
		return nil, nil
	}

	if room.Secret() != command.Token {
		h.logger.Error().
			Str("room_id", command.RoomID).
			Str("user_id", command.UserID).
			Msg("invalid secret key")

		return nil, ErrInvalidSecretKey
	}

	member.Authenticate()
	err = h.stge.Save(ctx, room)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("user authenticated")

	return nil, nil
}
