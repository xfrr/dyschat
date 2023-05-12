package rooms

import "context"

type Storage interface {
	Save(ctx context.Context, room *Room) (id string, err error)
	Get(ctx context.Context, id string) (*Room, error)
}
