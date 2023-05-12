package rooms

import (
	"time"

	"github.com/xfrr/dyschat/internal/events"
	"github.com/xfrr/dyschat/internal/pubsub"
)

type Status string

const (
	Opened Status = "opened"
	Closed Status = "closed"
)

type Room struct {
	events []pubsub.Event

	id        string
	name      string
	secretKey string
	members   []*Member
	ownerID   string
	status    Status
	createdAt time.Time
}

func NewRoom(id, name, ownerID, secretKey string) (*Room, error) {
	room := &Room{
		events: []pubsub.Event{
			&events.RoomCreatedEvent{
				RoomID:    id,
				RoomName:  name,
				SecretKey: secretKey,
				CreatedAt: time.Now().Unix(),
			},
		},

		id:        id,
		name:      name,
		secretKey: secretKey,
		ownerID:   id,
		status:    Opened,
		createdAt: time.Now(),
	}

	return room, room.validate()
}

func (r *Room) Events() []pubsub.Event {
	return r.events
}

func (r *Room) ID() string {
	return r.id
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) SecretKey() string {
	return r.secretKey
}

func (r *Room) Members() []*Member {
	return r.members
}

func (r *Room) Status() Status {
	return r.status
}

func (r *Room) CreatedAt() time.Time {
	return r.createdAt
}

func (r *Room) IsClosed() bool {
	return r.status == Closed
}

func (r *Room) IsOpened() bool {
	return r.status == Opened
}

func (r *Room) Open() {
	r.status = Opened
}

func (r *Room) Close() {
	r.status = Closed
}

func (r *Room) HasMember(id string) bool {
	for _, m := range r.members {
		if m.id == id {
			return true
		}
	}

	return false
}

func (r *Room) AddMember(username string) (*Member, error) {
	if r.HasMember(username) {
		return nil, ErrMemberAlreadyJoined
	}

	if r.isFull() {
		return nil, ErrRoomIsFull
	}

	member, err := newMember(username)
	if err != nil {
		return nil, err
	}

	r.members = append(r.members, member)

	r.events = append(r.events, &events.RoomMemberJoined{
		RoomID:   r.id,
		UserID:   username,
		JoinedAt: time.Now().Unix(),
	})

	return member, nil
}

func (r *Room) RemoveMember(id string) error {
	for i, m := range r.members {
		if m.id == id {
			r.members = append(r.members[:i], r.members[i+1:]...)
			return nil
		}
	}

	r.events = append(r.events, &events.RoomMemberLeft{
		RoomID:   r.id,
		MemberID: id,
	})

	return ErrMemberNotFound
}

func (r *Room) IsSecretKeyValid(secretKey string) bool {
	return r.secretKey == secretKey
}

func (r *Room) isFull() bool {
	return len(r.members) >= 2
}

func (r *Room) Restore(rm *RoomMemento) *Room {
	r.id = rm.ID
	r.name = rm.Name
	r.members = make([]*Member, 0, len(rm.Members))
	r.secretKey = rm.SecretKey
	r.status = rm.Status
	r.createdAt = rm.CreatedAt

	for _, m := range rm.Members {
		member := &Member{}
		r.members = append(r.members, member.Restore(m))
	}

	return r
}

func (r *Room) validate() error {
	if r.id == "" {
		return ErrEmptyRoomID
	}

	if r.name == "" {
		return ErrEmptyRoomName
	}

	if r.secretKey == "" {
		return ErrEmptyRoomSecretKey
	}

	return nil
}

type RoomMemento struct {
	ID        string
	Name      string
	SecretKey string
	Members   []*MemberMemento
	Status    Status
	CreatedAt time.Time
}
