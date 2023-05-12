package messages

import (
	"encoding/json"
	"time"

	"github.com/xfrr/dyschat/internal/events"
	"github.com/xfrr/dyschat/internal/pubsub"
)

type Metadata map[string]any

type MessageOptions func(*Message)

func WithMetadata(metadata Metadata) MessageOptions {
	return func(msg *Message) {
		msg.metadata = metadata
	}
}

func WithPublishedAt(publishedAt time.Time) MessageOptions {
	return func(msg *Message) {
		msg.publishedAt = publishedAt
	}
}

type Message struct {
	events []pubsub.Event

	id          string
	roomID      string
	userID      string
	text        string
	metadata    Metadata
	createdAt   time.Time
	publishedAt time.Time
}

func (m *Message) ID() string {
	return m.id
}

func (m *Message) RoomID() string {
	return m.roomID
}

func (m *Message) UserID() string {
	return m.userID
}

func (m *Message) Text() string {
	return m.text
}

func (m *Message) Metadata() Metadata {
	return m.metadata
}

func (m *Message) CreatedAt() time.Time {
	return m.createdAt
}

func (m *Message) PublishedAt() time.Time {
	return m.publishedAt
}

func (m *Message) MarshalJSON() ([]byte, error) {
	return json.Marshal(&MessageMemento{
		ID:          m.id,
		RoomID:      m.roomID,
		UserID:      m.userID,
		Text:        m.text,
		Metadata:    m.metadata,
		CreatedAt:   m.createdAt,
		PublishedAt: m.publishedAt,
	})
}

func (m *Message) UnmarshalJSON(data []byte) error {
	var msg MessageMemento
	if err := json.Unmarshal(data, &msg); err != nil {
		return err
	}

	m.id = msg.ID
	m.roomID = msg.RoomID
	m.userID = msg.UserID
	m.text = msg.Text
	m.metadata = msg.Metadata
	m.createdAt = msg.CreatedAt
	m.publishedAt = msg.PublishedAt

	return nil
}

func (m *Message) Restore(memento *MessageMemento) (*Message, error) {
	m.id = memento.ID
	m.roomID = memento.RoomID
	m.userID = memento.UserID
	m.text = memento.Text
	m.metadata = memento.Metadata
	m.createdAt = memento.CreatedAt
	m.publishedAt = memento.PublishedAt

	return m, m.validate()
}

func (m *Message) Events() []pubsub.Event {
	return m.events
}

func NewMessage(id string, roomID string, userID string, text string, opts ...MessageOptions) (*Message, error) {
	msg := &Message{
		id:        id,
		roomID:    roomID,
		userID:    userID,
		text:      text,
		createdAt: time.Now(),
	}

	for _, opt := range opts {
		opt(msg)
	}

	event := &events.MessageCreated{
		MessageID:   id,
		RoomID:      roomID,
		UserID:      userID,
		Text:        text,
		Metadata:    msg.metadata,
		CreatedAt:   msg.createdAt.Unix(),
		PublishedAt: msg.publishedAt.Unix(),
	}

	msg.events = append(msg.events, event)
	return msg, msg.validate()
}

func (m *Message) validate() error {
	if m.id == "" {
		return ErrInvalidMessageID
	}

	if m.roomID == "" {
		return ErrInvalidRoomID
	}

	if m.userID == "" {
		return ErrInvalidUserID
	}

	if m.text == "" {
		return ErrInvalidText
	}

	return nil
}

type MessageMemento struct {
	ID          string    `json:"id"`
	RoomID      string    `json:"room_id"`
	UserID      string    `json:"user_id"`
	Text        string    `json:"text"`
	Metadata    Metadata  `json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
}
