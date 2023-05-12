package events

import (
	"encoding/json"

	"github.com/xfrr/dyschat/internal/pubsub"
)

type RoomMemberJoined struct {
	RoomID   string `json:"room_id"`
	UserID   string `json:"user_id"`
	JoinedAt int64  `json:"joined_at"`
}

func (ev RoomMemberJoined) Subject() pubsub.Subject {
	return pubsub.Subject("rooms." + ev.RoomID + ".members." + ev.UserID + ".joined")
}

func (ev RoomMemberJoined) SubjectRegex() pubsub.Subject {
	return pubsub.Subject(`rooms\.([\d\w-_]+)\.members\.([\d\w-_]+)\.joined$`)
}

func (e *RoomMemberJoined) UnmarshalJSON(data []byte) error {
	type Alias RoomMemberJoined
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
