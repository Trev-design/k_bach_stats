package main

import (
	"log"
	"mailer_service/cmd/api"
)

func main() {
	if err := api.StartAndListen(); err != nil {
		log.Fatalf(err.Error())
	}

	log.Println("server is down")
}
