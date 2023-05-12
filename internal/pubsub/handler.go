package pubsub

import "context"

type Handlers map[Subject]Handler

type Handler interface {
	Handle(context.Context, string, []byte) error
}

func (h Handlers) Get(subject Subject) (Handler, bool) {
	// get handler subject by wildcard
	for s, h := range h {
		if s.Match(subject) {
			return h, true
		}
	}

	handler, ok := h[subject]
	return handler, ok
}
