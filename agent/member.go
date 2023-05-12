package agent

import (
	"context"
)

type MemberStatus int32

const (
	MemberStatusUnknown MemberStatus = iota
	MemberStatusConnected
	MemberStatusDisconnected
)

type Member struct {
	id              string
	status          MemberStatus
	isAuthenticated bool
	send            chan []byte
}

func newMember(id string) *Member {
	return &Member{
		id:              id,
		status:          MemberStatusUnknown,
		send:            make(chan []byte),
		isAuthenticated: false,
	}
}

func (c *Member) ID() string {
	return c.id
}

func (c *Member) Status() MemberStatus {
	return c.status
}

func (c *Member) Connect(send chan []byte) error {
	c.status = MemberStatusConnected
	c.send = send
	return nil
}

func (c *Member) Disconnect(ctx context.Context) error {
	c.status = MemberStatusDisconnected
	c.send = nil
	return nil
}

func (c *Member) IsConnected() bool {
	return c.status == MemberStatusConnected && c.send != nil
}

func (c *Member) IsDisconnected() bool {
	return c.status == MemberStatusDisconnected
}

func (c *Member) IsAuthenticated() bool {
	return c.isAuthenticated
}

func (c *Member) Authenticate() {
	c.isAuthenticated = true
}

func (c *Member) Send(ctx context.Context, message *Message) error {
	if !c.IsConnected() {
		return nil
	}

	msg, err := message.MarshalJson()
	if err != nil {
		return err
	}

	c.send <- msg
	return nil
}

func (c *Member) SendCh() chan []byte {
	return c.send
}

func (c *Member) Restore(
	id string,
	send chan []byte,
	isAuthenticated bool,
	status int32,
) *Member {
	c.id = id
	c.send = send
	c.isAuthenticated = isAuthenticated
	c.status = MemberStatus(status)
	return c
}
