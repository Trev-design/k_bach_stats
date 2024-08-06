package main

import (
	"broker-server/cmd/server"
	"log"
)

func main() {
	if err := server.StartAndListen(); err != nil {
		log.Fatalf("Failed to start the server: %v\n", err)
	}
}
