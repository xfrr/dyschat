package agent

import "errors"

var (
	ErrConnClosed            = errors.New("connection closed")
	ErrEmptyRoomID           = errors.New("empty room id")
	ErrEmptyRoomSecret       = errors.New("empty room secret")
	ErrMemberIsAlreadyInRoom = errors.New("member is already in room")
	ErrMemberNotFound        = errors.New("member not found")
	ErrRoomIsFull            = errors.New("room is full")
	ErrRoomNotFound          = errors.New("room not found")
	ErrUnauthenticated       = errors.New("unauthenticated")
	ErrUnknownMessageType    = errors.New("unknown message type")
)
