package pubsub

import "context"

type Consumer interface {
	Consume(ctx context.Context, topic string) error
}
