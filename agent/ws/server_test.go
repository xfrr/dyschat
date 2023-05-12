package ws_test

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func (s *Suite) newWebsocketConn(ctx context.Context, url string) (*websocket.Conn, error) {
	s.T().Logf("new websocket conn on url: %s", url)
	dial := websocket.Dialer{
		HandshakeTimeout: 3 * time.Second,
		Proxy:            http.ProxyFromEnvironment,
	}

	c, _, err := dial.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Suite) TestMaxConnections() {
	wg := sync.WaitGroup{}

	nconn := 130
	conns := 0
	for i := 0; i < nconn; i++ {
		wg.Add(1)
		go func() {
			conn, err := s.newWebsocketConn(context.Background(), "ws://localhost:50089/ws/health")
			s.NoError(err)

			defer func() {
				wg.Done()
			}()

			if conn != nil {
				conns++
			}
		}()
	}

	wg.Wait()
	s.T().Logf("conns: %d", conns)
	s.Equal(nconn, conns)
}
