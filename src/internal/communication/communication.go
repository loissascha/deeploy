package communication

import "encoding/json"

type MsgType string

const (
	MsgTypeInitiate MsgType = "initialize"
)

type WSMessage struct {
	MsgType string          `json:"msg_type"`
	Data    json.RawMessage `json:"data"`
}

type InitiateMessage struct {
	ID int `json:"id"`
}
