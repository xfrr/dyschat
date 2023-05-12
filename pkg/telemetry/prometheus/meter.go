package prometheus

import (
	"github.com/xfrr/dyschat/pkg/telemetry"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/sdk/metric"

	api "go.opentelemetry.io/otel/metric"
)

type Meter struct {
	meter api.Meter
}

func NewMeterProvider(svcName string, reader metric.Reader) *Meter {
	mp := metric.NewMeterProvider(metric.WithReader(reader))
	meter := mp.Meter(svcName)
	return &Meter{
		meter: meter,
	}
}

func (m *Meter) Int64Counter(name string, opts ...telemetry.CounterOption) (telemetry.Counter, error) {
	cfg := &telemetry.CounterOpts{}
	cfg.Apply(opts...)

	counter, err := m.meter.Int64Counter(
		name,
		instrument.WithDescription(cfg.Description()),
		instrument.WithUnit(cfg.Unit()),
	)
	if err != nil {
		return nil, err
	}

	return &Int64Counter{
		counter: counter,
	}, nil
}

func (m *Meter) Int64UpDownCounter(name string, opts ...telemetry.CounterOption) (telemetry.Counter, error) {
	cfg := &telemetry.CounterOpts{}
	cfg.Apply(opts...)

	counter, err := m.meter.Int64UpDownCounter(
		name,
		instrument.WithDescription(cfg.Description()),
		instrument.WithUnit(cfg.Unit()),
	)
	if err != nil {
		return nil, err
	}

	return &Int64UpDownCounter{
		counter: counter,
	}, nil
}
