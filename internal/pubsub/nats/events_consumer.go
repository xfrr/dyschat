package nats

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ pubsub.Consumer = (*EventsConsumer)(nil)

type EventsConsumer struct {
	nc *nats.Conn

	handlers pubsub.Handlers

	logger *zerolog.Logger

	cfg ConsumerOptions
}

func NewEventsConsumer(nc *nats.Conn, handlers pubsub.Handlers, logger *zerolog.Logger, opts ...ConsumerOption) (*EventsConsumer, error) {
	c := &EventsConsumer{
		nc:       nc,
		handlers: handlers,
		logger:   logger,
	}

	cfg := ConsumerOptions{}
	for _, opt := range opts {
		opt(&cfg)
	}

	return c, nil
}

func (c *EventsConsumer) Consume(ctx context.Context, topic string) error {
	c.logger.Debug().
		Str("topic", topic).
		Msg("subscribing to NATS topic")

	var err error
	_, err = c.nc.Subscribe(topic, func(msg *nats.Msg) {
		c.logger.Debug().
			Str("topic", topic).
			Str("subject", msg.Subject).
			Msg("received NATS message")

		handler, ok := c.handlers.Get(pubsub.Subject(msg.Subject))
		if !ok {
			if c.cfg.mustHaveHandler {
				c.logger.Debug().
					Str("topic", topic).
					Str("subject", msg.Subject).
					Msg("no handler found for NATS message")
			}
			return
		}

		err = handler.Handle(context.Background(), msg.Subject, msg.Data)
		if err != nil {
			c.logger.Debug().
				Str("topic", topic).
				Str("subject", msg.Subject).
				Err(err).
				Msg("failed to handle NATS message")
			return
		}
	})
	if err != nil {
		c.logger.Error().
			Str("topic", topic).
			Err(err).
			Msg("failed to subscribe to NATS topic")
		return err
	}

	return nil
}
