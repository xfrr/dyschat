package events

import (
	"encoding/json"

	"github.com/xfrr/dyschat/internal/pubsub"
)

type MessageCreated struct {
	MessageID   string         `json:"message_id"`
	RoomID      string         `json:"room_id"`
	UserID      string         `json:"user_id"`
	Text        string         `json:"text"`
	Metadata    map[string]any `json:"metadata"`
	CreatedAt   int64          `json:"created_at"`
	PublishedAt int64          `json:"published_at"`
}

func (ev MessageCreated) Subject() pubsub.Subject {
	return pubsub.Subject("rooms." + ev.RoomID + ".messages." + ev.MessageID + ".created")
}

func (ev MessageCreated) SubjectRegex() pubsub.Subject {
	return pubsub.Subject(`^rooms\.([\d\w-_]+)\.messages\.([\d\w-_]+)\.created$`)
}

func (e *MessageCreated) UnmarshalJSON(data []byte) error {
	type Alias MessageCreated
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
