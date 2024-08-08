package main

import (
	"log"
	"mailerservice/cmd/grpcserver"
)

func main() {
	if err := grpcserver.StartAndListen(); err != nil {
		log.Fatalf("huso")
	}
}
