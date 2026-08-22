package main

import (
	"local/deeploy/internal/server"
)

func main() {
	s := server.NewServer()

	go s.RunAgentHealthchecks()

	s.RunWS()
}
