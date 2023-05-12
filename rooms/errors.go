package rooms

import "errors"

var (
	ErrEmptyRoomID             = errors.New("empty room id")
	ErrEmptyRoomName           = errors.New("empty room name")
	ErrEmptyRoomSecretKey      = errors.New("empty room secret key")
	ErrRoomNotFound            = errors.New("room not found")
	ErrRoomInsufficientMembers = errors.New("room has insufficient members")

	ErrRoomIsClosed      = errors.New("room is closed")
	ErrRoomIsOpened      = errors.New("room is opened")
	ErrRoomAlreadyExists = errors.New("room already exists")
	ErrRoomIsFull        = errors.New("room is full")

	ErrEmptyMemberID       = errors.New("empty member id")
	ErrMemberNotFound      = errors.New("member not found")
	ErrMemberAlreadyJoined = errors.New("member already joined")
	ErrMemberNotJoined     = errors.New("member not joined")
	ErrMemberNotInRoom     = errors.New("member not in room")

	ErrUnauthorized = errors.New("unauthorized")
)
