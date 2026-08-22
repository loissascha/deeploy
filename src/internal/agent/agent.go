package agent

import (
	"context"
	"encoding/json"
	"local/deeploy/internal/communication"
	"log/slog"
	"time"

	"github.com/coder/websocket"
)

type Agent struct {
	id                    int
	conn                  *websocket.Conn
	lastHeartbeatReceived time.Time
}

func NewAgent(id int) *Agent {
	return &Agent{
		id:                    id,
		lastHeartbeatReceived: time.Now().UTC(),
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

func (a *Agent) RunReader(ctx context.Context) error {
	for {
		_, data, err := a.conn.Read(ctx)
		if err != nil {
			return err
		}

		var msg communication.WSMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			slog.Error("couldn't unmarshal received message.", "rawData", data, "err", err)
			continue
		}

		switch msg.MsgType {
		case communication.MsgTypeWelcome:
			slog.Info("server sent us a welcome message <3")
		case communication.MsgTypeError:
			slog.Error("server sent error message", "rawData", data)
		case communication.MsgTypeHeartbeat:
			slog.Info("server sent heartbeat.. sending one back <3")
			a.lastHeartbeatReceived = time.Now().UTC()
			err := a.WriteHeartbeatMessage(ctx)
			if err != nil {
				return err
			}
		default:
			slog.Error("unknown msg type:", "rawData", data)
		}
	}
}
