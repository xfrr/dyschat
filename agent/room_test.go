package agent_test

import (
	"context"
	"time"

	"github.com/xfrr/dyschat/agent"
)

type roomArgs struct {
	id        string
	secretKey string
	createdAt time.Time
}

func (suite *AgentStorageTestSuite) TestNewRoom() {
	tests := []struct {
		name string

		args roomArgs
		err  error
	}{
		{
			name: "valid room",
			args: roomArgs{
				id:        "room-id",
				secretKey: "secret-key",
				createdAt: time.Now(),
			},
		},
		{
			name: "empty room id",
			err:  agent.ErrEmptyRoomID,
			args: roomArgs{
				id:        "",
				secretKey: "secret-key",
				createdAt: time.Now(),
			},
		},
		{
			name: "empty secret key",
			err:  agent.ErrEmptyRoomSecret,
			args: roomArgs{
				id:        "test-id",
				secretKey: "",
				createdAt: time.Now(),
			},
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			room, err := agent.NewRoom(tt.args.id, tt.args.secretKey, tt.args.createdAt)
			if tt.err != nil {
				suite.EqualError(err, tt.err.Error())
				suite.Nil(room)
			} else {
				suite.NoError(err)
				suite.NotNil(room)
				suite.Equal(tt.args.id, room.ID())
				suite.Equal(tt.args.secretKey, room.Secret())
				suite.Equal(tt.args.createdAt, room.CreatedAt())
				suite.Empty(room.Members())
			}
		})
	}
}

func (suite *AgentStorageTestSuite) TestMembers() {
	room, err := agent.NewRoom("test-id", "test-secret", time.Now())
	suite.NoError(err)

	m1, err := room.AddMember("m1")
	suite.NoError(err)
	suite.Equal(m1.ID(), "m1")
	suite.Equal(m1.Status(), agent.MemberStatusUnknown)
	suite.NotZero(m1.SendCh())
	suite.False(m1.IsAuthenticated())

	m2, err := room.AddMember("m2")
	suite.NoError(err)
	suite.Equal(m2.ID(), "m2")
	suite.Equal(m2.Status(), agent.MemberStatusUnknown)
	suite.NotZero(m2.SendCh())
	suite.False(m2.IsAuthenticated())

	suite.Equal(2, len(room.Members()))

	// test already in room error
	_, err = room.AddMember("m2")
	suite.Error(err)

	// get member 1
	m1, err = room.GetMember("m1")
	suite.NoError(err)

	// test member not found error
	_, err = room.GetMember("m3")
	suite.Error(err)

	send := make(chan []byte, 1)
	err = m1.Connect(send)
	suite.NoError(err)
	err = m2.Connect(send)
	suite.NoError(err)

	room.Broadcast(context.Background(), m1.ID(), agent.NewEventMessage("message-id", nil))

	suite.Assert().Eventually(func() bool {
		return len(m2.SendCh()) == 1
	}, time.Second, 10*time.Millisecond)

	// test remove members
	room.RemoveMember(*m1)
	room.RemoveMember(*m2)

	suite.Empty(room.Members())
}
