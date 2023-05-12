package grpc

import (
	context "context"
	"errors"

	"github.com/rs/zerolog"

	"github.com/xfrr/dyschat/rooms"
	"github.com/xfrr/dyschat/rooms/commands"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	iauth "github.com/xfrr/dyschat/internal/auth"
	icommands "github.com/xfrr/dyschat/internal/commands"
	proto "github.com/xfrr/dyschat/proto/rooms/v1"
)

var _ proto.RoomsServiceServer = (*RoomsServer)(nil)

type RoomsServer struct {
	proto.UnsafeRoomsServiceServer

	cmdbus  *icommands.Bus
	querier rooms.Querier
	logger  *zerolog.Logger

	authInterceptor *iauth.AuthGRPCInterceptor
}

func NewRoomsServer(
	cmdbus *icommands.Bus,
	querier rooms.Querier,
	authInterceptor *iauth.AuthGRPCInterceptor,
	logger *zerolog.Logger,
) *RoomsServer {
	return &RoomsServer{
		querier:         querier,
		cmdbus:          cmdbus,
		logger:          logger,
		authInterceptor: authInterceptor,
	}
}

func (srv *RoomsServer) GetRoomMessages(ctx context.Context, req *proto.GetRoomMessagesRequest) (*proto.Messages, error) {
	// messages, err := srv.msgquerier.GetMessages(ctx, req.GetRoomId(), req.GetStartTimestamp(), req.GetEndTimestamp())
	// if err != nil {
	// 	return nil, toGrpcError(err)
	// }

	// pbmessages := make([]*proto.Message, 0, len(messages))
	// for _, m := range messages {
	// 	pbmessages = append(pbmessages, &proto.Message{
	// 		Id:        m.ID(),
	// 		RoomId:    m.RoomID(),
	// 		MemberId:  m.UserID(),
	// 		Content:   m.Content(),
	// 		Timestamp: m.CreatedAt().Unix(),
	// 	})
	// }

	// return &proto.Messages{
	// 	Messages: pbmessages,
	// }, nil
	return nil, errors.New("not implemented")
}

func (s *RoomsServer) JoinMember(ctx context.Context, req *proto.JoinMemberRequest) (*proto.Empty, error) {
	cmd := &commands.JoinMemberCommand{
		RoomID:    req.GetRoomId(),
		UserID:    req.GetMemberId(),
		SecretKey: req.GetSecretKey(),
	}

	_, err := s.cmdbus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, toGrpcError(err)
	}

	return &proto.Empty{}, nil
}

func (s *RoomsServer) LeaveMember(ctx context.Context, req *proto.LeaveMemberRequest) (*proto.Empty, error) {
	cmd := &commands.LeaveMemberCommand{
		RoomID: req.GetRoomId(),
		UserID: req.GetMemberId(),
	}

	_, err := s.cmdbus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, toGrpcError(err)
	}

	return &proto.Empty{}, nil
}

func (rs *RoomsServer) CreateRoom(ctx context.Context, req *proto.CreateRoomRequest) (*proto.ID, error) {
	uid, err := iauth.UserIDFromContext(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, errors.New("unauthenticated").Error())
	}

	cmd := &commands.CreateRoomCommand{
		ID:        req.GetId(),
		Name:      req.GetName(),
		SecretKey: req.GetSecretKey(),
		OwnerID:   uid,
	}

	r, err := rs.cmdbus.Dispatch(ctx, cmd)
	if err != nil {
		return nil, toGrpcError(err)
	}

	reply := r.(*commands.CreateRoomReply)
	return &proto.ID{Id: reply.RoomID, SecretKey: reply.SecretKey}, nil
}

func (rs *RoomsServer) GetRoom(ctx context.Context, req *proto.GetRoomRequest) (*proto.Room, error) {
	room, err := rs.querier.GetRoom(ctx, req.GetRoomId())
	if err != nil {
		return nil, toGrpcError(err)
	}

	members := make([]*proto.Member, 0, len(room.Members()))
	for _, m := range room.Members() {
		members = append(members, &proto.Member{
			Id:     m.ID(),
			Status: proto.MemberStatus(proto.MemberStatus_value[string(m.Status())]),
		})
	}

	return &proto.Room{
		Id:        room.ID(),
		Name:      room.Name(),
		SecretKey: room.SecretKey(),
		Members:   members,
		Status:    proto.RoomStatus(proto.RoomStatus_value[string(room.Status())]),
	}, nil
}

func (rs *RoomsServer) GetRooms(ctx context.Context, req *proto.GetRoomsRequest) (*proto.Rooms, error) {
	rooms, err := rs.querier.ListRooms(ctx)
	if err != nil {
		return nil, toGrpcError(err)
	}

	roomsProto := make([]*proto.Room, 0, len(rooms))
	for _, room := range rooms {
		protoRoom := &proto.Room{
			Id:        room.ID(),
			Name:      room.Name(),
			SecretKey: room.SecretKey(),
			Status:    proto.RoomStatus(proto.RoomStatus_value[string(room.Status())]),
		}

		for _, m := range room.Members() {
			protoRoom.Members = append(protoRoom.Members, &proto.Member{
				Id:     m.ID(),
				Status: proto.MemberStatus(proto.MemberStatus_value[string(m.Status())]),
			})
		}

		roomsProto = append(roomsProto, protoRoom)
	}

	return &proto.Rooms{
		Rooms: roomsProto,
	}, nil
}

func toGrpcError(err error) error {
	switch err {
	case rooms.ErrRoomNotFound:
		return status.Error(codes.NotFound, err.Error())
	case rooms.ErrRoomAlreadyExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case rooms.ErrRoomIsFull:
		return status.Error(codes.ResourceExhausted, err.Error())
	case rooms.ErrRoomIsClosed:
		return status.Error(codes.FailedPrecondition, err.Error())
	case rooms.ErrUnauthorized, rooms.ErrMemberNotInRoom, rooms.ErrMemberNotJoined:
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return err
	}
}
