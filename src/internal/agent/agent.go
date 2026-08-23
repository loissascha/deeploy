package agent

import (
	"context"
	"encoding/json"
	"local/deeploy/internal/communication"
	"local/deeploy/internal/settings"
	"log/slog"
	"time"

	"github.com/coder/websocket"
)

type Agent struct {
	settings              *settings.AgentSettings
	conn                  *websocket.Conn
	lastHeartbeatReceived time.Time
}

func NewAgent(s *settings.AgentSettings) *Agent {
	return &Agent{
		settings:              s,
		lastHeartbeatReceived: time.Now().UTC(),
	}
}

func (a *Agent) RunConn(ctx context.Context) error {
	conn, _, err := websocket.Dial(ctx, a.settings.ControllerHostWS, nil)
	if err != nil {
		return err
	}
	a.conn = conn
	err = a.writeInitializeMessage(ctx)
	if err != nil {
		return err
	}
	err = a.runReader(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (a *Agent) Close() {
	if a.conn != nil {
		a.conn.CloseNow()
	}
}

func (a *Agent) writeInitializeMessage(ctx context.Context) error {
	raw, err := json.Marshal(communication.InitiateMessage{
		AgentID: a.settings.AgentID,
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

func (a *Agent) writeHeartbeatMessage(ctx context.Context) error {
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

func (a *Agent) handleWelcomeMessage(msg communication.WSMessage) error {
	slog.Info("server sent us a welcome message <3")

	var w communication.WelcomeMessage
	err := json.Unmarshal(msg.Data, &w)
	if err != nil {
		return err
	}

	if a.settings.AgentID != w.YourID {
		a.settings.AgentID = w.YourID
		err = a.settings.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Agent) runReader(ctx context.Context) error {
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
			err := a.handleWelcomeMessage(msg)
			if err != nil {
				slog.Error("error in communication message", "err", err)
			}
		case communication.MsgTypeError:
			slog.Error("server sent error message", "rawData", data)
		case communication.MsgTypeHeartbeat:
			slog.Info("server sent heartbeat.. sending one back <3")
			a.lastHeartbeatReceived = time.Now().UTC()
			err := a.writeHeartbeatMessage(ctx)
			if err != nil {
				return err
			}
		default:
			slog.Error("unknown msg type:", "rawData", data)
		}
	}
}
