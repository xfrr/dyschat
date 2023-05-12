package nats

import "github.com/nats-io/nats.go"

func NewStreams(streams ...nats.StreamConfig) []nats.StreamConfig {
	return streams
}

func NewJetStreamContext(nc *nats.Conn, streams []nats.StreamConfig) (nats.JetStreamContext, error) {
	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, err
	}

	for _, stream := range streams {
		if _, err := js.StreamInfo(stream.Name); err == nil {
			continue
		}

		_, err := js.AddStream(&stream)
		if err != nil {
			return nil, err
		}
	}

	return js, nil
}
