package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/dyschat/internal/pubsub"
)

type StreamPublisher struct {
	nc     nats.JetStreamContext
	stream string
}

func NewStreamPublisher(nc nats.JetStreamContext, stream string) *StreamPublisher {
	return &StreamPublisher{
		nc:     nc,
		stream: stream,
	}
}

func (p *StreamPublisher) Publish(ctx context.Context, events []pubsub.Event) error {
	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}

		_, err = p.nc.Publish(string(event.Subject()), data)
		if err != nil {
			return err
		}

	}

	return nil
}
