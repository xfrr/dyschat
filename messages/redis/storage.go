package redis

import (
	"context"
	"encoding"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/xfrr/dyschat/messages"
)

var _ messages.Storage = (*Storage)(nil)

var _ encoding.BinaryMarshaler = (*messageDTO)(nil)

type Storage struct {
	rdb *redis.Client
}

type messageDTO struct {
	ID          string      `redis:"id"`
	RoomID      string      `redis:"room_id"`
	UserID      string      `redis:"user_id"`
	Text        string      `redis:"text"`
	Metadata    interface{} `redis:"metadata"`
	CreatedAt   int64       `redis:"created_at"`
	PublishedAt int64       `redis:"published_at"`
}

func (t *messageDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *messageDTO) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}

func NewStorage(rdb *redis.Client) messages.Storage {
	return &Storage{
		rdb: rdb,
	}
}

func (stg *Storage) Save(ctx context.Context, message *messages.Message) (err error) {
	r := &messageDTO{
		ID:          message.ID(),
		RoomID:      message.RoomID(),
		UserID:      message.UserID(),
		Text:        message.Text(),
		Metadata:    message.Metadata(),
		CreatedAt:   message.CreatedAt().Unix(),
		PublishedAt: message.PublishedAt().Unix(),
	}

	err = stg.rdb.Set(ctx, message.ID(), r, 0).Err()
	if err != nil {
		return err
	}

	err = stg.rdb.HSet(ctx, "messages", message.ID(), r).Err()
	if err != nil {
		return err
	}
	return nil
}
