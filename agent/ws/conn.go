package ws

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/xfrr/dyschat/agent"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 3) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Conn is an middleman between the websocket connection and the hub.
type Conn struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered room of outbound messages.
	send chan []byte

	// The room id this connection is in.
	roomID string

	// The member id this connection is for.
	memberID string

	headers http.Header

	authenticated bool
}

func (c *Conn) Disconnect(ctx context.Context) error {
	close(c.send)
	return c.ws.Close()
}

func (c *Conn) Send(ctx context.Context, msg *agent.Message) error {
	if c.send == nil {
		return agent.ErrConnClosed
	}

	payload, err := msg.MarshalJson()
	if err != nil {
		return err
	}

	c.send <- payload
	return nil
}

func (c *Conn) start(ctx context.Context, received chan *message, unregister chan *Conn) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if !c.authenticated {
		go c.checkAuth(ctx, unregister)
	}

	go c.writePump(ctx)
	c.readPump(ctx, received)
}

// write writes a message with the given message type and payload.
func (c *Conn) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

func (conn *Conn) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-conn.send:
			if !ok {
				// The hub closed the room.
				conn.write(websocket.CloseMessage, []byte{})
				return
			}

			conn.ws.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := conn.ws.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(conn.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-conn.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := conn.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (conn *Conn) readPump(ctx context.Context, received chan *message) error {
	conn.ws.SetReadLimit(maxMessageSize)
	conn.ws.SetReadDeadline(time.Now().Add(pongWait))
	conn.ws.SetPongHandler(func(string) error {
		conn.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			_, msg, err := conn.ws.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					return err
				}
				return nil
			}

			msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))

			received <- &message{
				from: conn,
				data: msg,
			}
		}
	}
}

func (conn *Conn) checkAuth(ctx context.Context, unregister chan *Conn) {
	// check if the connection is authenticated every 5 seconds and timeout at 30 seconds
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	timeout := time.After(90 * time.Second)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if conn.memberID != "" {
				return
			}
			conn.sendError(ctx, errors.New("you are not authenticated, please, authenticate before 90 seconds or you will be disconnected"), false)
		case <-timeout:
			unregister <- conn
			return
		}
	}
}

type errorMessage struct {
	Err string `json:"error"`
}

func (c *Conn) sendError(ctx context.Context, err error, close bool) {
	errmsg := toErrorMessage(err)
	payload, _ := json.Marshal(errmsg)
	c.write(websocket.TextMessage, payload)
	if close {
		c.Disconnect(ctx)
	}
}

func toErrorMessage(err error) errorMessage {
	switch err {
	case agent.ErrUnauthenticated, agent.ErrMemberNotFound:
		return errorMessage{Err: agent.ErrUnauthenticated.Error()}
	default:
		return errorMessage{Err: err.Error()}
	}
}
