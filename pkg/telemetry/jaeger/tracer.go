package jaeger

import (
	"context"

	"github.com/xfrr/dyschat/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

var _ telemetry.Tracer = (*Tracer)(nil)

type Tracer struct {
	serviceName string

	provider *tracesdk.TracerProvider
}

func (t *Tracer) StartSpan(ctx context.Context, name string, opts ...telemetry.SpanOption) (context.Context, telemetry.Span) {
	span := &Span{
		name:        name,
		serviceName: t.serviceName,
	}

	for _, opt := range opts {
		opt(span)
	}

	return span.start(ctx, t.provider)
}

func (t Tracer) Shutdown(ctx context.Context) error {
	return t.provider.Shutdown(ctx)
}

func NewTracerProvider(url, serviceName, envMode string) (*Tracer, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			attribute.String("environment", envMode),
		)),
	)

	return &Tracer{
		serviceName: serviceName,
		provider:    tp,
	}, nil
}
