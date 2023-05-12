package commands

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/internal/commands"
	"github.com/xfrr/dyschat/internal/idn"
	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ commands.CommandHandler = (*ConnectCommandHandler)(nil)

type ConnectCommandHandler struct {
	idp       idn.Provider
	stge      agent.RoomStorage
	publisher pubsub.Publisher

	logger *zerolog.Logger
}

func NewConnectCommandHandler(idp idn.Provider, stge agent.RoomStorage, pub pubsub.Publisher, logger *zerolog.Logger) *ConnectCommandHandler {
	return &ConnectCommandHandler{
		idp:       idp,
		stge:      stge,
		publisher: pub,
		logger:    logger,
	}
}

func (h *ConnectCommandHandler) Handle(ctx context.Context, cmd commands.Command) (commands.Reply, error) {
	command, ok := cmd.(*ConnectCommand)
	if !ok {
		return nil, commands.NewErrInvalidCommandType(cmd.Type())
	}

	h.logger.Debug().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("connecting user to room")

	if err := command.validate(); err != nil {
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

	if !member.IsAuthenticated() {
		return nil, agent.ErrUnauthenticated
	}

	if !member.IsConnected() {
		err = member.Connect(command.Send)
		if err != nil {
			return nil, err
		}
	}

	err = h.stge.Save(ctx, room)
	if err != nil {
		return nil, err
	}

	h.logger.Info().
		Str("room_id", command.RoomID).
		Str("user_id", command.UserID).
		Msg("user connected to room")

	// TODO: publish MemberConnectedEvent

	return nil, nil
}
