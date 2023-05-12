package agent

import (
	"context"
	"time"
)

const (
	maxRoomSize = 10
)

type Room struct {
	id        string
	secret    string
	members   []*Member
	createdAt time.Time
}

func NewRoom(id, secretKey string, createdAt time.Time) (*Room, error) {
	room := &Room{
		id:        id,
		secret:    secretKey,
		members:   []*Member{},
		createdAt: createdAt,
	}

	if err := room.validate(); err != nil {
		return nil, err
	}

	return room, nil
}

func (r Room) ID() string {
	return r.id
}

func (r Room) Secret() string {
	return r.secret
}

func (r Room) Members() []*Member {
	return r.members
}

func (r Room) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Room) Broadcast(ctx context.Context, senderID string, messages ...*Message) (err error) {
	for _, m := range messages {
		for _, member := range r.members {
			if senderID == member.id {
				continue
			}

			err = member.Send(ctx, m)
			if err != nil {
				continue
			}
		}
	}

	return err
}

func (r *Room) AddMember(id string) (*Member, error) {
	if r.IsFull() {
		return nil, ErrRoomIsFull
	}
	if r.IsMember(id) {
		return nil, ErrMemberIsAlreadyInRoom
	}
	member := newMember(id)
	r.members = append(r.members, member)
	return member, nil
}

func (r *Room) RemoveMember(member Member) {
	for i, m := range r.members {
		if m.id == member.id {
			r.members = append(r.members[:i], r.members[i+1:]...)
			return
		}
	}
}

func (r Room) GetMember(id string) (*Member, error) {
	for _, m := range r.members {
		if m.id == id {
			return m, nil
		}
	}
	return nil, ErrMemberNotFound
}

func (r Room) IsMember(id string) bool {
	for _, m := range r.members {
		if m.id == id {
			return true
		}
	}
	return false
}

func (r Room) IsEmpty() bool {
	return len(r.members) == 0
}

func (r Room) IsFull() bool {
	return len(r.members) == maxRoomSize
}

func (r *Room) Restore(
	id string,
	secret string,
	members []*Member,
	createdAt time.Time,
) *Room {
	r.id = id
	r.members = members
	r.secret = secret
	r.createdAt = createdAt
	return r
}

func (r *Room) validate() error {
	if len(r.id) == 0 {
		return ErrEmptyRoomID
	}

	if len(r.secret) == 0 {
		return ErrEmptyRoomSecret
	}

	return nil
}
