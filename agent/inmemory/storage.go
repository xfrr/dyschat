package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/xfrr/dyschat/agent"
)

var _ agent.RoomStorage = (*RoomStorage)(nil)

type RoomStorage struct {
	rooms []*roomDTO

	mu sync.Mutex
}

type roomDTO struct {
	id        string
	secret    string
	members   []*memberDTO
	createdAt int64
}

type memberDTO struct {
	id              string
	send            chan []byte
	isAuthenticated bool
	status          agent.MemberStatus
}

func NewRoomStorage() *RoomStorage {
	return &RoomStorage{
		rooms: []*roomDTO{},
		mu:    sync.Mutex{},
	}
}

func (s *RoomStorage) Save(ctx context.Context, room *agent.Room) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, r := range s.rooms {
		if r.id == room.ID() {
			s.rooms[i] = s.toDTO(room)
			return nil
		}
	}

	s.rooms = append(s.rooms, s.toDTO(room))
	return nil
}

func (s *RoomStorage) Get(ctx context.Context, id string) (room *agent.Room, err error) {
	for _, r := range s.rooms {
		if r.id == id {
			return toRoom(r), nil
		}
	}

	return nil, agent.ErrRoomNotFound
}

func (s *RoomStorage) Delete(ctx context.Context, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, room := range s.rooms {
		if room.id == id {
			s.rooms = append(s.rooms[:i], s.rooms[i+1:]...)
		}
	}
}

func (s *RoomStorage) toDTO(room *agent.Room) *roomDTO {
	members := []*memberDTO{}
	for _, member := range room.Members() {
		members = append(members, &memberDTO{
			id:              member.ID(),
			send:            member.SendCh(),
			isAuthenticated: member.IsAuthenticated(),
			status:          member.Status(),
		})
	}

	return &roomDTO{
		id:        room.ID(),
		secret:    room.Secret(),
		members:   members,
		createdAt: room.CreatedAt().Unix(),
	}
}

func toRoom(dto *roomDTO) *agent.Room {
	members := []*agent.Member{}
	for _, m := range dto.members {
		member := agent.Member{}
		member.Restore(
			m.id,
			m.send,
			m.isAuthenticated,
			int32(m.status),
		)

		members = append(members, &member)
	}

	room := agent.Room{}
	room.Restore(
		dto.id,
		dto.secret,
		members,
		time.Unix(dto.createdAt, 0),
	)

	return &room
}
