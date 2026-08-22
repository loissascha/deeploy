package server

import (
	"context"

	"github.com/coder/websocket"
)

type ServerAgent struct {
	id   int
	conn *websocket.Conn
}

func NewServerAgent(id int, conn *websocket.Conn) *ServerAgent {
	return &ServerAgent{
		id:   id,
		conn: conn,
	}
}

func (a *ServerAgent) SendWelcomeMessage(ctx context.Context) error {
	err := a.conn.Write(ctx, websocket.MessageText, []byte("welcome"))
	if err != nil {
		return err
	}
	return nil
}
