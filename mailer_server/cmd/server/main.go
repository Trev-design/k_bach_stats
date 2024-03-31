package main

import (
	"fmt"
	"mailer-server/cmd/api"
)

func main() {
	fmt.Println("starting mail server")
	api.NewApiServer(":8080").Run()
	fmt.Println("ciao")
}
