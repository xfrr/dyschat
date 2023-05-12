package telemetry

import "context"

type Counter interface {
	Add(ctx context.Context, value int64, opts ...CounterAttributes)
}

type CounterOpts struct {
	description string
	unit        string
}

type CounterOption func(*CounterOpts)

type CounterAttributes []CounterAttribute

type CounterAttribute struct {
	Key   string
	Value string
}

func WithDescription(description string) CounterOption {
	return func(o *CounterOpts) {
		o.description = description
	}
}

func WithUnit(unit string) CounterOption {
	return func(o *CounterOpts) {
		o.unit = unit
	}
}

func (o *CounterOpts) Description() string {
	return o.description
}

func (o *CounterOpts) Unit() string {
	return o.unit
}

func (o *CounterOpts) Apply(opts ...CounterOption) {
	for _, opt := range opts {
		opt(o)
	}
}
