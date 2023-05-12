package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ pubsub.Consumer = (*PersistentStreamConsumer)(nil)

type PersistentStreamConsumer struct {
	nc nats.JetStreamContext

	handlers pubsub.Handlers

	logger *zerolog.Logger

	cfg ConsumerOptions
}

func NewPersistentStreamConsumer(nc nats.JetStreamContext, handlers pubsub.Handlers, logger *zerolog.Logger, opts ...ConsumerOption) (*PersistentStreamConsumer, error) {
	cfg := ConsumerOptions{
		mustHaveHandler: true,
		maxWait:         5 * time.Second,
	}
	for _, opt := range opts {
		opt(&cfg)
	}

	c := &PersistentStreamConsumer{
		nc:       nc,
		handlers: handlers,
		logger:   logger,
		cfg:      cfg,
	}

	return c, nil
}

func (c *PersistentStreamConsumer) Consume(ctx context.Context, subject string) error {
	c.logger.Debug().
		Str("subject", subject).
		Msg("starting to consume nats events from stream")

	if c.cfg.bindStream != "" && c.cfg.durable != "" {
		c.nc.AddConsumer(subject, &nats.ConsumerConfig{
			Durable:   c.cfg.durable,
			AckPolicy: nats.AckExplicitPolicy,
		})
	}

	sub, err := c.nc.PullSubscribe(subject, c.cfg.durable, nats.BindStream(c.cfg.bindStream))
	if err != nil {
		return err
	}

L:
	for {
		select {
		case <-ctx.Done():
			break L
		default:
			msgs, err := sub.Fetch(10, nats.MaxWait(c.cfg.maxWait))
			if err == nats.ErrTimeout {
				if c.cfg.closeOnTimeout {
					break L
				}
			} else if err != nil {
				return err
			}

			if len(msgs) == 0 {
				continue
			}

			for _, msg := range msgs {
				msg.Ack()
				handler, ok := c.handlers.Get(pubsub.Subject(msg.Subject))
				if !ok {
					if c.cfg.mustHaveHandler {
						c.logger.Error().
							Str("subject", msg.Subject).
							Msg("no handler found for nats event")
					}
					continue
				}

				err = handler.Handle(context.Background(), msg.Subject, msg.Data)
				if err != nil {
					c.logger.Debug().
						Err(err).
						Str("subject", msg.Subject).
						Msg("failed to handle NATS message")
					continue
				}
			}
		}
	}

	return c.drainSubscription(sub, subject)
}

func (c *PersistentStreamConsumer) drainSubscription(sub *nats.Subscription, subject string) error {
	c.logger.Debug().
		Str("subject", subject).
		Msg("draining nats subscription")

	return sub.Drain()
}
