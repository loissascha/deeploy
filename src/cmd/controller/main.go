package main

import (
	"local/deeploy/internal/server"
)

func main() {
	s := server.NewServer()

	s.RunWS()
}
