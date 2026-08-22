package server

import (
	"context"
	"encoding/json"
	"local/deeploy/internal/communication"
	"log"
	"log/slog"
	"time"

	"github.com/coder/websocket"
)

type ServerAgent struct {
	id                    int
	conn                  *websocket.Conn
	lastHeartbeatReceived time.Time
}

func NewServerAgent(id int, conn *websocket.Conn) *ServerAgent {
	return &ServerAgent{
		id:                    id,
		conn:                  conn,
		lastHeartbeatReceived: time.Now().UTC(),
	}
}

func (a *ServerAgent) SendWelcomeMessage(ctx context.Context) error {
	raw, err := json.Marshal(communication.WelcomeMessage{
		Msg: "welcome",
	})
	if err != nil {
		return err
	}

	msg := communication.WSMessage{
		MsgType: communication.MsgTypeWelcome,
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

func (a *ServerAgent) SendErrorMessage(ctx context.Context, str string) error {
	raw, err := json.Marshal(communication.ErrorMessage{
		Error: str,
	})
	if err != nil {
		return err
	}

	msg := communication.WSMessage{
		MsgType: communication.MsgTypeError,
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

func (a *ServerAgent) SendHeartbeatMessage(ctx context.Context) error {
	slog.Info("sending heartbeat to agent", "id", a.id)
	raw, err := json.Marshal(communication.HeartbeatMessage{
		Msg: "are you still alive?", // TODO: maybe some sort of secret/identifier for the server?
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

func (a *ServerAgent) RunReader(ctx context.Context) {
	defer a.conn.CloseNow()

	for {
		_, data, err := a.conn.Read(ctx)
		if err != nil {
			log.Println("read:", err)
			return
		}

		var msg communication.WSMessage
		err = json.Unmarshal(data, &msg)
		if err != nil {
			log.Println("unmarshal:", err)
			return
		}

		switch msg.MsgType {
		case communication.MsgTypeInitiate:
			slog.Error("why is there a duplicated initiate message")
		case communication.MsgTypeHeartbeat:
			slog.Info("heartbeat received on server")
			a.lastHeartbeatReceived = time.Now().UTC()
		default:
			slog.Warn("unknown msg type:", "msgtype", msg.MsgType, "data", msg.Data)
		}
	}
}
