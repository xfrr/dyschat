package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog"
	"github.com/xfrr/dyschat/agent"
	"github.com/xfrr/dyschat/agent/commands"
	"github.com/xfrr/dyschat/pkg/telemetry"

	icommands "github.com/xfrr/dyschat/internal/commands"
)

var _ agent.Server = (*Server)(nil)

type message struct {
	from *Conn
	data []byte
}

type Server struct {
	addr string
	path string

	received   chan *message
	register   chan *Conn
	unregister chan *Conn

	logger *zerolog.Logger

	bus *icommands.Bus

	connections map[*Conn]struct{}

	meter telemetry.Meter
}

type ServerOption func(*Server)

func WithAddr(addr string) ServerOption {
	return func(srv *Server) {
		srv.addr = addr
	}
}

func WithPath(path string) ServerOption {
	return func(srv *Server) {
		srv.path = path
	}
}

func WithMeter(meter telemetry.Meter) ServerOption {
	return func(srv *Server) {
		srv.meter = meter
	}
}

func NewServer(bus *icommands.Bus, logger *zerolog.Logger, opts ...ServerOption) *Server {
	srv := &Server{
		addr:        ":9000",
		path:        "/rooms/{room_id}/messages",
		received:    make(chan *message),
		register:    make(chan *Conn),
		unregister:  make(chan *Conn),
		logger:      logger,
		bus:         bus,
		connections: make(map[*Conn]struct{}),
		meter:       &telemetry.NoopMeter{},
	}

	for _, opt := range opts {
		opt(srv)
	}

	return srv
}

func (srv *Server) Start(ctx context.Context) error {
	go srv.start(ctx)

	srv.logger.Info().
		Str("addr", srv.addr).
		Msg("starting websocket server")

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route(srv.path, func(r chi.Router) {
		r.Get("/{room_id}/live", srv.handler)
		r.Get("/health", srv.healthHandler)
	})

	return http.ListenAndServe(srv.addr, r)
}

func (srv *Server) start(ctx context.Context) {
	defer srv.stop(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case conn, ok := <-srv.register:
			if !ok {
				srv.logger.Error().Msg("register channel closed")
				continue
			}
			if err := srv.registerConn(context.Background(), conn); err != nil {
				srv.logger.Error().
					Err(err).
					Str("room_id", conn.roomID).
					Str("user_id", conn.memberID).
					Msg("error registering connection")
				conn.sendError(context.Background(), err, true)
			}
		case conn, ok := <-srv.unregister:
			if !ok {
				srv.logger.Error().Msg("unregister channel closed")
				continue
			}

			srv.logger.Debug().
				Str("room_id", conn.roomID).
				Str("addr", conn.ws.RemoteAddr().String()).
				Msg("unregister connection")

			if _, ok := srv.connections[conn]; ok {
				delete(srv.connections, conn)
				err := conn.Disconnect(context.Background())
				if err != nil {
					srv.logger.Debug().Err(err).Msg("stop connection")
				}
			}
		case msg, ok := <-srv.received:
			if !ok {
				srv.logger.Error().Msg("received channel closed")
				continue
			}

			err := srv.handleMessage(context.Background(), msg)
			if err != nil {
				srv.logger.Error().Err(err).Msg("handle message")
				msg.from.sendError(context.Background(), err, false)
			}
		}
	}
}

func (srv *Server) stop(ctx context.Context) {
	for conn := range srv.connections {
		err := conn.Disconnect(ctx)
		if err != nil {
			srv.logger.Debug().
				Str("room_id", conn.roomID).
				Str("addr", conn.ws.RemoteAddr().String()).
				Err(err).
				Msg("closed connection")
		}
	}

	srv.logger.Info().Msg("websocket server stopped")
}

func (srv *Server) handler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		srv.logger.Error().Err(err).Msg("upgrade")
		return
	}

	send := make(chan []byte, 256)
	conn := &Conn{
		send:    send,
		ws:      ws,
		headers: r.Header,
	}

	conn.roomID = chi.URLParam(r, "room_id")
	if conn.roomID == "" {
		srv.logger.Error().Msg("required url param room_id is empty")
		conn.sendError(context.Background(), errors.New("required url param room_id is empty"), true)
		return
	}

	conn.memberID = r.Header.Get("X-User-ID")

	srv.logger.Info().
		Str("room_id", conn.roomID).
		Str("user_id", conn.memberID).
		Str("addr", r.RemoteAddr).
		Interface("headers", r.Header).
		Msg("new connection")

	srv.register <- conn
}

func (srv *Server) handleMessage(ctx context.Context, msg *message) error {
	type tmp struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}

	var t tmp
	err := json.Unmarshal(msg.data, &t)
	if err != nil {
		srv.logger.Error().Err(err).Msg("unmarshal message")
		return err
	}

	switch t.Type {
	case string(commands.ConnectCommandType):
		return srv.dispatchConnectCommand(ctx, msg.from, t.Payload)
	case string(commands.PublishMessageCommandType):
		return srv.dispatchPublishMessageCommand(ctx, msg.from, t.Payload)
	default:
		return agent.ErrUnknownMessageType
	}
}

func (srv *Server) dispatchConnectCommand(ctx context.Context, from *Conn, payload json.RawMessage) error {
	cmd := commands.ConnectCommand{
		RoomID: from.roomID,
		Send:   from.send,
	}

	err := json.Unmarshal(payload, &cmd)
	if err != nil {
		return err
	}

	_, err = srv.bus.Dispatch(ctx, &commands.AuthenticateCommand{
		RoomID: cmd.RoomID,
		UserID: cmd.UserID,
		Token:  cmd.SecretKey,
	})
	if err != nil {
		return err
	}

	_, err = srv.bus.Dispatch(ctx, &cmd)
	if err != nil {
		return err
	}

	from.memberID = cmd.UserID
	from.authenticated = true
	return nil
}

func (srv *Server) dispatchPublishMessageCommand(ctx context.Context, from *Conn, payload json.RawMessage) error {
	var cmd commands.PublishMessageCommand
	err := json.Unmarshal(payload, &cmd)
	if err != nil {
		return err
	}

	cmd.RoomID = from.roomID
	cmd.UserID = from.memberID

	_, err = srv.bus.Dispatch(ctx, &cmd)
	if err != nil {
		return err
	}

	return nil
}

func (srv *Server) registerConn(ctx context.Context, conn *Conn) error {
	token := conn.headers.Get("X-Auth-Token")
	if token != "" && conn.memberID != "" {
		_, err := srv.bus.Dispatch(ctx, &commands.AuthenticateCommand{
			RoomID: conn.roomID,
			UserID: conn.memberID,
			Token:  token,
		})
		if err != nil {
			return err
		}

		conn.authenticated = true

		_, err = srv.bus.Dispatch(ctx, &commands.ConnectCommand{
			RoomID: conn.roomID,
			UserID: conn.memberID,
			Send:   conn.send,
		})
		if err != nil {
			return err
		}
	}

	srv.connections[conn] = struct{}{}
	go conn.start(ctx, srv.received, srv.unregister)

	counter, _ := srv.meter.Int64UpDownCounter("connections")
	counter.Add(ctx, 1)
	return nil
}

func (srv *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
		ReadBufferSize:    1024,
		WriteBufferSize:   1024,
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		srv.logger.Error().Err(err).Msg("upgrade")
		return
	}

	err = ws.WriteMessage(websocket.TextMessage, []byte("ok"))
	if err != nil {
		srv.logger.Error().Err(err).Msg("write message")
		return
	}
}
