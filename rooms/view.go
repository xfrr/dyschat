package rooms

import "context"

type View interface {
	Get(ctx context.Context, id string) (*Room, error)
	List(ctx context.Context, opts ...ListOption) ([]*Room, error)
}

type ListOption func(*ListOptions)

type ListOptions struct {
	ids      []string
	orderBy  string
	asceding bool
}

func (o *ListOptions) IDs() []string {
	return o.ids
}

func (o *ListOptions) OrderBy() string {
	return o.orderBy
}

func (o *ListOptions) Asceding() bool {
	return o.asceding
}

func WithOrderBy(orderBy string) ListOption {
	return func(o *ListOptions) {
		o.orderBy = orderBy
	}
}

func WithAsceding(asceding bool) ListOption {
	return func(o *ListOptions) {
		o.asceding = asceding
	}
}

func WithIDs(ids ...string) ListOption {
	return func(o *ListOptions) {
		o.ids = ids
	}
}

func (r *ListOptions) Apply(opts ...ListOption) {
	for _, opt := range opts {
		opt(r)
	}
}
