package messages

import "context"

type Storage interface {
	Save(ctx context.Context, message *Message) error
}
