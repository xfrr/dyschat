package telemetry

import "context"

type Meter interface {
	Int64Counter(name string, opts ...CounterOption) (Counter, error)
	Int64UpDownCounter(name string, opts ...CounterOption) (Counter, error)
}

type NoopMeter struct{}

type NoopCounter struct{}

func (m *NoopMeter) Int64Counter(name string, opts ...CounterOption) (Counter, error) {
	return &NoopCounter{}, nil
}

func (m *NoopMeter) Int64UpDownCounter(name string, opts ...CounterOption) (Counter, error) {
	return &NoopCounter{}, nil
}

func (c *NoopCounter) Add(ctx context.Context, value int64, opts ...CounterAttributes) {}
