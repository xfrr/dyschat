package nats

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	"github.com/xfrr/dyschat/internal/pubsub"
)

type EventsPublisher struct {
	nc *nats.Conn
}

func NewEventsPublisher(nc *nats.Conn) *EventsPublisher {
	return &EventsPublisher{
		nc: nc,
	}
}

func (p *EventsPublisher) Publish(ctx context.Context, events []pubsub.Event) error {
	for _, event := range events {
		data, err := json.Marshal(event)
		if err != nil {
			return err
		}

		err = p.nc.Publish(string(event.Subject()), data)
		if err != nil {
			return err
		}
	}

	return nil
}
