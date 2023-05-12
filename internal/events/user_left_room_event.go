package events

import (
	"encoding/json"

	"github.com/xfrr/dyschat/internal/pubsub"
)

type RoomMemberLeft struct {
	RoomID   string `json:"room_id"`
	MemberID string `json:"member_id"`
	LeftAt   int64  `json:"left_at"`
}

func (ev RoomMemberLeft) Subject() pubsub.Subject {
	return pubsub.Subject("rooms." + ev.RoomID + ".members." + ev.MemberID + ".left")
}

func (ev RoomMemberLeft) SubjectRegex() pubsub.Subject {
	return pubsub.Subject(`^rooms\.([\d\w-_]+)\.members\.([\d\w-_]+)\.left$`)
}

func (e *RoomMemberLeft) UnmarshalJSON(data []byte) error {
	type Alias RoomMemberLeft
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
