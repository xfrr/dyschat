package ws_test

import (
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
	"github.com/xfrr/dyschat/agent/ws"
	"github.com/xfrr/dyschat/internal/commands"
)

type Suite struct {
	suite.Suite

	srv     *ws.Server
	stopSrv func()
}

func (s *Suite) SetupSuite() {
	bus := commands.NewBus()
	logger := zerolog.Nop()
	s.srv = ws.NewServer(
		bus,
		&logger,
		ws.WithAddr(":50089"),
		ws.WithPath("/ws"),
	)
	ctx, cancel := context.WithCancel(context.Background())
	go s.srv.Start(ctx)

	s.stopSrv = cancel
}

func (s *Suite) TearDownSuite() {
	s.stopSrv()
}

func (s *Suite) SetupTest() {
	s.T().Log("setup test")
}

func (s *Suite) TearDownTest() {

}

func BenchmarkSuite(b *testing.B) {
	s := new(Suite)
	s.SetT(&testing.T{})
	s.SetupSuite()
	b.ResetTimer()

	count := 0
	for i := 0; i < b.N; i++ {
		s.SetupTest()
		b.StartTimer()

		count++
		s.newWebsocketConn(context.Background(), "ws://localhost:50089/ws/health")

		b.StopTimer()
		s.TearDownTest()
	}

	b.Logf("count: %d", count)
	s.TearDownSuite()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
