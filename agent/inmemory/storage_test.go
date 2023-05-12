package inmemory_test

import (
	"context"
	"time"

	"github.com/xfrr/dyschat/agent"
)

func (suite *InMemoryStorageTestSuite) TestSaveNewRoom() {
	room, err := agent.NewRoom(
		"room-id",
		"room-secret",
		time.Now(),
	)
	suite.NoError(err)

	err = suite.sut.Save(context.Background(), room)
	suite.NoError(err)

	found, err := suite.sut.Get(context.Background(), "room-id")
	suite.NoError(err)
	suite.NotNil(room)

	suite.assertEqualRooms(room, found)
}

func (suite *InMemoryStorageTestSuite) TestSaveExistingRoom() {
	room, err := agent.NewRoom(
		"room-id",
		"room-secret",
		time.Now(),
	)
	suite.NoError(err)

	err = suite.sut.Save(context.Background(), room)
	suite.NoError(err)

	member, err := room.AddMember("test-member")
	suite.NoError(err)

	ch := make(chan []byte)
	err = member.Connect(ch)
	suite.NoError(err)

	suite.NoError(err)

	err = suite.sut.Save(context.Background(), room)
	suite.NoError(err)

	found, err := suite.sut.Get(context.Background(), "room-id")
	suite.NoError(err)

	suite.assertEqualRooms(room, found)
}

func (suite *InMemoryStorageTestSuite) TestDeleteRoom() {
	room, err := agent.NewRoom(
		"room-id",
		"room-secret",
		time.Now(),
	)

	suite.NoError(err)

	err = suite.sut.Save(context.Background(), room)
	suite.NoError(err)

	suite.sut.Delete(context.Background(), "room-id")
	suite.NoError(err)

	found, err := suite.sut.Get(context.Background(), "room-id")
	suite.Error(err)
	suite.Nil(found)
}

func (suite *InMemoryStorageTestSuite) assertEqualRooms(expected, actual *agent.Room) {
	suite.Equal(expected.ID(), actual.ID())
	suite.Equal(expected.Secret(), actual.Secret())
	suite.Equal(expected.CreatedAt().Format(time.RFC3339), actual.CreatedAt().Format(time.RFC3339))
	for i, m := range expected.Members() {
		suite.Equal(m.ID(), actual.Members()[i].ID())
		suite.Equal(m.Status(), actual.Members()[i].Status())
		suite.Equal(m.SendCh(), actual.Members()[i].SendCh())
		suite.Equal(m.IsAuthenticated(), actual.Members()[i].IsAuthenticated())
		suite.Equal(m.IsAuthenticated(), actual.Members()[i].IsAuthenticated())
	}
}
