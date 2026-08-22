package communication

import "encoding/json"

type MsgType string

const (
	MsgTypeInitiate  MsgType = "initialize"
	MsgTypeHeartbeat MsgType = "heartbeat"
)

type WSMessage struct {
	MsgType MsgType         `json:"msg_type"`
	Data    json.RawMessage `json:"data,omitempty"`
}

type InitiateMessage struct {
	ID int `json:"id"`
}

type HeartbeatMessage struct {
	Msg string `json:"msg"`
}
