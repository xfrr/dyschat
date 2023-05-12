package telemetry

import (
	"context"
)

type SpanOption func(Span)

type Span interface {
	Error(err error)
	End() error
}

type Tracer interface {
	StartSpan(ctx context.Context, name string, opts ...SpanOption) (newCtx context.Context, span Span)
	Shutdown(context.Context) error
}
