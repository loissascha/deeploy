package agent

import (
	"context"
	"encoding/json"
	"local/deeploy/internal/communication"
	"log"
	"log/slog"

	"github.com/coder/websocket"
)

type Agent struct {
	id   int
	conn *websocket.Conn
}

func NewAgent(id int) *Agent {
	return &Agent{
		id: id,
	}
}

func (a *Agent) Dial(ctx context.Context) error {
	conn, _, err := websocket.Dial(ctx, "ws://localhost:42066/ws", nil)
	if err != nil {
		return err
	}
	a.conn = conn
	return nil
}

func (a *Agent) Close() {
	if a.conn != nil {
		a.conn.CloseNow()
	}
}

func (a *Agent) WriteInitializeMessage(ctx context.Context) error {
	raw, err := json.Marshal(communication.InitiateMessage{
		ID: a.id,
	})
	if err != nil {
		return err
	}

	msg := communication.WSMessage{
		MsgType: communication.MsgTypeInitiate,
		Data:    raw,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = a.conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return err
	}

	return nil
}

func (a *Agent) WriteHeartbeatMessage(ctx context.Context) error {
	slog.Info("sending heartbeat to server")
	raw, err := json.Marshal(communication.HeartbeatMessage{
		Msg: "I'm still there.",
	})
	if err != nil {
		return err
	}

	msg := communication.WSMessage{
		MsgType: communication.MsgTypeHeartbeat,
		Data:    raw,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	err = a.conn.Write(ctx, websocket.MessageText, data)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) WriteTest(ctx context.Context) error {
	err := a.conn.Write(
		ctx,
		websocket.MessageText,
		[]byte("hello controller"),
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) RunReader(ctx context.Context) error {
	for {
		_, data, err := a.conn.Read(ctx)
		if err != nil {
			return err
		}
		log.Printf("received: %s\n", data)
	}
}
