package events

import (
	"context"

	"github.com/xfrr/dyschat/internal/pubsub"
	"github.com/xfrr/dyschat/pkg/telemetry"
)

type EventHandlerTelemetryMiddleware struct {
	tracer telemetry.Tracer

	handler pubsub.Handler
}

func NewEventHandlerTelemetryMiddleware(tracer telemetry.Tracer, handler pubsub.Handler) *EventHandlerTelemetryMiddleware {
	return &EventHandlerTelemetryMiddleware{
		tracer:  tracer,
		handler: handler,
	}
}

func (c *EventHandlerTelemetryMiddleware) Handle(ctx context.Context, subject string, payload []byte) error {
	ctx, span := c.tracer.StartSpan(ctx, "event."+subject)
	defer span.End()

	err := c.handler.Handle(ctx, subject, payload)
	if err != nil {
		span.Error(err)
	}

	return err
}
