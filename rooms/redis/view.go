package redis

import (
	"context"
	"encoding"
	"encoding/json"

	"github.com/redis/go-redis/v9"
	"github.com/xfrr/dyschat/rooms"
)

var _ rooms.View = (*View)(nil)

var _ encoding.BinaryMarshaler = (*roomDTO)(nil)
var _ encoding.BinaryMarshaler = (*memberDTO)(nil)

type View struct {
	rdb *redis.Client
}

func (stg *View) Get(ctx context.Context, id string) (*rooms.Room, error) {
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

func (stg *View) List(ctx context.Context, opts ...rooms.ListOption) ([]*rooms.Room, error) {
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

func NewView(rdb *redis.Client) *View {
	return &View{
		rdb: rdb,
	}
}
