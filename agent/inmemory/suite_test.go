package inmemory_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/xfrr/dyschat/agent/inmemory"
)

type InMemoryStorageTestSuite struct {
	suite.Suite
	sut *inmemory.RoomStorage
}

func (suite *InMemoryStorageTestSuite) SetupTest() {
	suite.sut = inmemory.NewRoomStorage()
}

func TestInMemoryStorageTestSuite(t *testing.T) {
	suite.Run(t, new(InMemoryStorageTestSuite))
}
