package main

import (
	"fmt"
	"mailer-server/cmd/api"
)

func main() {
	api.NewApiServer(":8080").Run()
	fmt.Println("ciao")
}
