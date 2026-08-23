package communication

import "encoding/json"

type MsgType string

const (
	MsgTypeInitiate  MsgType = "initialize"
	MsgTypeWelcome   MsgType = "welcome"
	MsgTypeHeartbeat MsgType = "heartbeat"
	MsgTypeError     MsgType = "error"
)

type WSMessage struct {
	MsgType MsgType         `json:"msg_type"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type InitiateMessage struct {
	AgentID int `json:"id"`
}

type WelcomeMessage struct {
	YourID int `json:"id"`
}

type HeartbeatMessage struct {
	Msg string `json:"msg,omitempty"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
