package main

import (
	"context"
	"local/deeploy/internal/agent"
	"log/slog"
	"time"
)

func main() {
	a := agent.NewAgent(1)
	ctx := context.Background()

	go func() {
		for {
			if err := a.RunConn(ctx); err != nil {
				slog.Error("error running the server conn", "err", err)
			}
			slog.Info("retrying in 10 seconds...")
			time.Sleep(10 * time.Second)
		}
	}()

	for {
	}
}
