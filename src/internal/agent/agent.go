package agent

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

type Agent struct {
	id int
}

func NewAgent(id int) *Agent {
	return &Agent{
		id: id,
	}
}

func (a *Agent) Run(ctx context.Context) error {
	conn, _, err := websocket.Dial(ctx, "ws://localhost:42066/ws", nil)
	if err != nil {
		return err
	}
	defer conn.CloseNow()

	err = conn.Write(
		ctx,
		websocket.MessageText,
		[]byte("hello controller"),
	)
	if err != nil {
		return err
	}

	_, data, err := conn.Read(ctx)
	if err != nil {
		return err
	}

	log.Printf("received: %s\n", data)

	return nil
}
