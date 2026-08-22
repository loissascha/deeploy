package main

import (
	"context"
	"local/deeploy/internal/agent"
	"log/slog"
	"time"
)

// TODO:
// has jobs (how are they defined? in settings? in whatever? per cli?)
// after receiving the welcome message -> register each job
// each job can have a run

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
