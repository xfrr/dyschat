package agent_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type AgentStorageTestSuite struct {
	suite.Suite
}

func (suite *AgentStorageTestSuite) SetupTest() {
}

func TestAgentStorageTestSuite(t *testing.T) {
	suite.Run(t, new(AgentStorageTestSuite))
}
