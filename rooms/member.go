package rooms

import "time"

type MemberStatus int32

const (
	MemberStatusUnknown MemberStatus = iota
	MemberStatusOnline
	MemberStatusOffline
)

type Member struct {
	id              string
	lastMessageID   int64
	lastMessageAt   time.Time
	lastConnectedAt time.Time
	status          MemberStatus
}

type MemberMemento struct {
	ID              string
	LastMessageID   int64
	LastMessageAt   time.Time
	LastConnectedAt time.Time
	Status          MemberStatus
}

func newMember(id string) (*Member, error) {
	member := &Member{
		id:     id,
		status: MemberStatusUnknown,
	}

	return member, member.validate()
}

func (m *Member) ID() string {
	return m.id
}

func (m *Member) LastMessageID() int64 {
	return m.lastMessageID
}

func (m *Member) SetLastMessageID(id int64) {
	m.lastMessageID = id
}

func (m *Member) LastMessageAt() time.Time {
	return m.lastMessageAt
}

func (m *Member) SetLastMessageAt(t time.Time) {
	m.lastMessageAt = t
}

func (m *Member) LastConnectedAt() time.Time {
	return m.lastConnectedAt
}

func (m *Member) Status() MemberStatus {
	return m.status
}

func (m *Member) SetOnline() {
	m.status = MemberStatusOnline
	m.lastConnectedAt = time.Now()
}

func (m *Member) SetOffline() {
	m.status = MemberStatusOffline
	m.lastConnectedAt = time.Now()
}

func (m *Member) Restore(mm *MemberMemento) *Member {
	m.id = mm.ID
	m.lastMessageID = mm.LastMessageID
	m.lastMessageAt = mm.LastMessageAt
	m.lastConnectedAt = mm.LastConnectedAt
	m.status = mm.Status
	return m
}

func (m *Member) validate() error {
	if m.id == "" {
		return ErrEmptyMemberID
	}

	return nil
}
