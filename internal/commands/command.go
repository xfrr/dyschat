package commands

import "context"

type Type string

type Reply any

type Command interface {
	Type() Type
}

type CommandHandler interface {
	Handle(context.Context, Command) (Reply, error)
}

type BusOption func(*Bus)

func WithCommand(command Command, handler CommandHandler) BusOption {
	return func(opts *Bus) {
		opts.handlers[command.Type()] = handler
	}
}

type Bus struct {
	handlers map[Type]CommandHandler
}

func NewBus(options ...BusOption) *Bus {
	bus := &Bus{
		handlers: make(map[Type]CommandHandler),
	}
	for _, opt := range options {
		opt(bus)
	}

	return bus
}

func (b *Bus) Dispatch(ctx context.Context, cmd Command) (Reply, error) {
	handler, ok := b.handlers[cmd.Type()]
	if !ok {
		return nil, ErrCommandHandlerNotFound
	}

	return handler.Handle(ctx, cmd)
}

func (b *Bus) Register(command Command, handler CommandHandler) {
	b.handlers[command.Type()] = handler
}

func (b *Bus) Unregister(ctype Type) error {
	delete(b.handlers, ctype)
	return nil
}
