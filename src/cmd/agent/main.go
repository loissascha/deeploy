package main

import (
	"context"
	"local/deeploy/internal/agent"
	"log"
)

func main() {
	a := agent.NewAgent(1)
	ctx := context.Background()

	if err := a.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
