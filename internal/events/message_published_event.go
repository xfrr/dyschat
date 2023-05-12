package events

import (
	"encoding/json"

	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ pubsub.Event = (*MessagePublishedEvent)(nil)

type MessagePublishedEvent struct {
	MessageID   string         `json:"message_id"`
	RoomID      string         `json:"room_id"`
	UserID      string         `json:"user_id"`
	Text        string         `json:"text"`
	Metadata    map[string]any `json:"metadata"`
	PublishedAt int64          `json:"published_at"`
}

func (ev MessagePublishedEvent) Subject() pubsub.Subject {
	return pubsub.Subject("rooms." + ev.RoomID + ".messages." + ev.MessageID + ".published")
}

func (ev MessagePublishedEvent) SubjectRegex() pubsub.Subject {
	return pubsub.Subject(`rooms\.([\d\w-_]+)\.messages\.([\d\w-_]+)\.published$`)
}

func (e *MessagePublishedEvent) UnmarshalJSON(data []byte) error {
	type Alias MessagePublishedEvent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}
