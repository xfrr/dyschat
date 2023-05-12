package agent

import (
	"encoding/json"
	"time"
)

type MessageType string

const (
	MessageTypeEvent MessageType = "event"
	MessageTypeError MessageType = "error"
)

type message struct {
	ID        string          `json:"id"`
	Type      MessageType     `json:"type"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

type Message struct {
	id        string
	mtype     MessageType
	payload   []byte
	createdAt time.Time
}

func (m Message) ID() string {
	return m.id
}

func (m Message) Type() MessageType {
	return m.mtype
}

func (m Message) IsEvent() bool {
	return m.mtype == MessageTypeEvent
}

func (m Message) IsError() bool {
	return m.mtype == MessageTypeError
}

func (m Message) Payload() []byte {
	return m.payload
}

func (m Message) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Message) MarshalJson() ([]byte, error) {
	return json.Marshal(&message{
		ID:        m.id,
		Type:      m.mtype,
		Payload:   m.payload,
		CreatedAt: m.createdAt,
	})
}

func (m *Message) UnmarshalJson(data []byte) error {
	msg := &message{}
	if err := json.Unmarshal(data, msg); err != nil {
		return err
	}
	m.id = msg.ID
	m.mtype = msg.Type
	m.payload = msg.Payload
	m.createdAt = msg.CreatedAt
	return nil
}

func NewEventMessage(id string, payload []byte) *Message {
	return &Message{
		id:        id,
		mtype:     MessageTypeEvent,
		payload:   payload,
		createdAt: time.Now(),
	}
}

func NewErrorMessage(id string, payload []byte) *Message {
	return &Message{
		id:        id,
		mtype:     MessageTypeError,
		payload:   payload,
		createdAt: time.Now(),
	}
}
