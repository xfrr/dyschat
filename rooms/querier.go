package rooms

import "context"

type Querier interface {
	GetRoom(ctx context.Context, id string) (*Room, error)
	ListRooms(ctx context.Context, opts ...ListOption) ([]*Room, error)
}

type querier struct {
	view View
}

func (q *querier) GetRoom(ctx context.Context, id string) (*Room, error) {
	return q.view.Get(ctx, id)
}

func (q *querier) ListRooms(ctx context.Context, opts ...ListOption) ([]*Room, error) {
	return q.view.List(ctx, opts...)
}

func NewQuerier(view View) *querier {
	return &querier{
		view: view,
	}
}
