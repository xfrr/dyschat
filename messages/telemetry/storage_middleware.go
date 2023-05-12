package telemetry

import (
	"context"

	"github.com/xfrr/dyschat/messages"
	"github.com/xfrr/dyschat/pkg/telemetry"
)

var _ messages.Storage = (*StorageTelemetryMiddleware)(nil)

type StorageTelemetryMiddleware struct {
	tracer telemetry.Tracer
	meter  telemetry.Meter

	stge messages.Storage
}

func NewStorageTelemetryMiddleware(tracer telemetry.Tracer, meter telemetry.Meter, stge messages.Storage) *StorageTelemetryMiddleware {
	return &StorageTelemetryMiddleware{
		tracer: tracer,
		meter:  meter,
		stge:   stge,
	}
}

func (s *StorageTelemetryMiddleware) Save(ctx context.Context, message *messages.Message) error {
	ctx, span := s.tracer.StartSpan(ctx, "storage.messages.save")
	defer span.End()

	err := s.stge.Save(ctx, message)
	if err != nil {
		return err
	}

	counter, _ := s.meter.Int64Counter("messages.saved")
	counter.Add(ctx, 1)
	return nil
}
