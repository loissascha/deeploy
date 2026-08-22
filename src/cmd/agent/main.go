package main

import (
	"context"
	"local/deeploy/internal/agent"
	"log"
)

func main() {
	a := agent.NewAgent(1)
	ctx := context.Background()

	if err := a.Dial(ctx); err != nil {
		log.Fatal(err)
	}

	if err := a.WriteTest(ctx); err != nil {
		log.Fatal(err)
	}

	if err := a.RunReader(ctx); err != nil {
		log.Fatal(err)
	}
}
