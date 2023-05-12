package commands

import (
	"context"

	"github.com/xfrr/dyschat/pkg/telemetry"
)

type CommandHandlerMiddleware struct {
	tracer telemetry.Tracer

	handler CommandHandler
}

func NewCommandHandlerMiddleware(tracer telemetry.Tracer, handler CommandHandler) *CommandHandlerMiddleware {
	return &CommandHandlerMiddleware{
		tracer:  tracer,
		handler: handler,
	}
}

func (c *CommandHandlerMiddleware) Handle(ctx context.Context, cmd Command) (Reply, error) {
	ctx, span := c.tracer.StartSpan(ctx, "command."+string(cmd.Type()))
	defer span.End()

	res, err := c.handler.Handle(ctx, cmd)
	if err != nil {
		span.Error(err)
	}

	return res, err
}
