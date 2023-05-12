package nats

import "time"

type ConsumerOptions struct {
	// MustHaveHandler specifies whether the consumer must have a handler for
	// every message it receives.
	mustHaveHandler bool

	maxWait time.Duration

	closeOnTimeout bool

	bindStream string

	durable string
}

type ConsumerOption func(*ConsumerOptions)

// WithMustHaveHandler specifies whether the consumer must have a handler for
// every message it receives.
func WithMustHaveHandler(must bool) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.mustHaveHandler = must
	}
}

// WithMaxWait specifies the maximum amount of time the consumer will wait for
// a message to arrive before returning.
func WithMaxWait(d time.Duration) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.maxWait = d
	}
}

func WithCloseOnDone() ConsumerOption {
	return func(o *ConsumerOptions) {
		o.closeOnTimeout = true
	}
}

func WithBindStream(stream string) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.bindStream = stream
	}
}

func WithDurable(durable string) ConsumerOption {
	return func(o *ConsumerOptions) {
		o.durable = durable
	}
}
