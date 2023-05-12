package events

import (
	"encoding/json"

	"github.com/xfrr/dyschat/internal/pubsub"
)

type MemberConnected struct {
	RoomID      string `json:"room_id"`
	MemberID    string `json:"member_id"`
	ConnectedAt int64  `json:"connected_at"`
}

func (ev MemberConnected) Subject() pubsub.Subject {
	return pubsub.Subject("rooms." + ev.RoomID + ".members." + ev.MemberID + ".connected")
}

func (ev MemberConnected) SubjectRegex() pubsub.Subject {
	return pubsub.Subject(`^rooms\.([\d\w-_]+)\.members\.([\d\w-_]+)\.connected$`)
}

func (e *MemberConnected) UnmarshalJSON(data []byte) error {
	type Alias MemberConnected
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
