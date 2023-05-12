package agent

import "context"

//go:generate moq -out mock/room_storage.go -pkg mock . RoomStorage:RoomStorageMock
type RoomStorage interface {
	Save(ctx context.Context, room *Room) error
	Get(ctx context.Context, roomID string) (*Room, error)
	Delete(ctx context.Context, roomID string)
}
