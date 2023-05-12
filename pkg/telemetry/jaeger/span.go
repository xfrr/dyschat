package jaeger

import (
	"context"

	"github.com/xfrr/dyschat/pkg/telemetry"
	"go.opentelemetry.io/otel/trace"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

var _ telemetry.Span = (*Span)(nil)

type Span struct {
	serviceName string
	name        string

	span trace.Span
}

func (s Span) Error(err error) {
	s.span.RecordError(err)
}

func (s Span) End() error {
	s.span.End()
	return nil
}

func (s *Span) start(ctx context.Context, provider *tracesdk.TracerProvider) (context.Context, Span) {
	newCtx, tspan := provider.Tracer(s.serviceName).Start(ctx, s.name)
	s.span = tspan
	return newCtx, *s
}
