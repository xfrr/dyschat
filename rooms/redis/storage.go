package redis

import (
	"context"
	"encoding"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xfrr/dyschat/rooms"
)

var _ rooms.Storage = (*Storage)(nil)

var _ encoding.BinaryMarshaler = (*roomDTO)(nil)
var _ encoding.BinaryMarshaler = (*memberDTO)(nil)

type Storage struct {
	rdb *redis.Client
}

type roomDTO struct {
	ID        string      `redis:"id"`
	Name      string      `redis:"name"`
	SecretKey string      `redis:"secret_key"`
	Members   []memberDTO `redis:"members"`
	Status    string      `redis:"status"`
	CreatedAt int64       `redis:"created_at"`
}

type memberDTO struct {
	ID              string `redis:"id"`
	LastMessageSeq  int64  `redis:"last_message_id"`
	LastMessageAt   int64  `redis:"last_message_at"`
	LastConnectedAt int64  `redis:"last_connected"`
	Status          int32  `redis:"status"`
}

func (t *roomDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *roomDTO) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}

func (t *memberDTO) MarshalBinary() ([]byte, error) {
	return json.Marshal(t)
}

func (t *memberDTO) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}

	return nil
}

func (stg *Storage) Get(ctx context.Context, id string) (*rooms.Room, error) {
	raw, err := stg.rdb.Get(ctx, id).Result()
	if err != nil {
		return nil, rooms.ErrRoomNotFound
	}

	var r roomDTO
	err = json.Unmarshal([]byte(raw), &r)
	if err != nil {
		return nil, err
	}

	return toRoom(r), nil
}

func (stg *Storage) List(ctx context.Context, opts ...rooms.ListOption) ([]*rooms.Room, error) {
	var results []*rooms.Room

	cmd := stg.rdb.HGetAll(ctx, "rooms")
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	for _, v := range cmd.Val() {
		var r roomDTO
		err := json.Unmarshal([]byte(v), &r)
		if err != nil {
			return nil, err
		}

		results = append(results, toRoom(r))
	}

	return results, nil
}

func (stg *Storage) Save(ctx context.Context, room *rooms.Room) (id string, err error) {
	var members []memberDTO
	for _, m := range room.Members() {
		members = append(members, memberDTO{
			ID:              m.ID(),
			LastMessageSeq:  m.LastMessageID(),
			LastMessageAt:   m.LastMessageAt().Unix(),
			LastConnectedAt: m.LastConnectedAt().Unix(),
			Status:          int32(m.Status()),
		})
	}

	r := &roomDTO{
		ID:        room.ID(),
		Name:      room.Name(),
		SecretKey: room.SecretKey(),
		Members:   members,
		Status:    string(room.Status()),
		CreatedAt: room.CreatedAt().Unix(),
	}

	err = stg.rdb.Set(ctx, room.ID(), r, 0).Err()
	if err != nil {
		return room.ID(), err
	}

	err = stg.rdb.HSet(ctx, "rooms", room.ID(), r).Err()
	if err != nil {
		return room.ID(), err
	}

	return room.ID(), nil
}

func NewStorage(rdb *redis.Client) *Storage {
	return &Storage{
		rdb: rdb,
	}
}

func toRoom(dto roomDTO) *rooms.Room {
	m := []*rooms.MemberMemento{}
	for _, mDTO := range dto.Members {
		m = append(m, &rooms.MemberMemento{
			ID:              mDTO.ID,
			LastMessageID:   mDTO.LastMessageSeq,
			LastMessageAt:   time.Unix(mDTO.LastMessageAt, 0),
			LastConnectedAt: time.Unix(mDTO.LastConnectedAt, 0),
			Status:          rooms.MemberStatus(mDTO.Status),
		})
	}

	room := &rooms.Room{}
	return room.Restore(&rooms.RoomMemento{
		ID:        dto.ID,
		Name:      dto.Name,
		SecretKey: dto.SecretKey,
		Members:   m,
		Status:    rooms.Status(dto.Status),
		CreatedAt: time.Unix(dto.CreatedAt, 0),
	})
}
