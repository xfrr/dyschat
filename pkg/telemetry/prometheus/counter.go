package prometheus

import (
	"context"

	"github.com/xfrr/dyschat/pkg/telemetry"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/instrument"
)

type Int64Counter struct {
	counter instrument.Int64Counter
}

func (c *Int64Counter) Add(ctx context.Context, value int64, opts ...telemetry.CounterAttributes) {
	c.counter.Add(ctx, value, c.attributes(opts...)...)
}

func (c *Int64Counter) attributes(opts ...telemetry.CounterAttributes) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(opts))
	for _, opt := range opts {
		for _, attr := range opt {
			attrs = append(attrs, attribute.String(attr.Key, attr.Value))
		}
	}

	return attrs
}
