package events

import (
	"encoding/json"

	"github.com/xfrr/dyschat/internal/pubsub"
)

var _ pubsub.Event = (*RoomCreatedEvent)(nil)

type RoomCreatedEvent struct {
	RoomID    string `json:"room_id"`
	RoomName  string `json:"room_name"`
	SecretKey string `json:"secret_key"`
	CreatedAt int64  `json:"created_at"`
}

func (rce RoomCreatedEvent) Subject() pubsub.Subject {
	return pubsub.Subject("rooms." + rce.RoomID + ".events.created")
}

func (rce RoomCreatedEvent) SubjectRegex() pubsub.Subject {
	return pubsub.Subject(`rooms\.([\d\w-_]+)\.events.created`)
}

func (rce *RoomCreatedEvent) UnmarshalJSON(data []byte) error {
	type Alias RoomCreatedEvent
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(rce),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	return nil
}
