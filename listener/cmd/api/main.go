package main

import (
	"listener/cmd/server"
	"log"
)

func main() {
	if err := server.StartAndListen(); err != nil {
		log.Fatalf("%v\n", err)
	}
}
