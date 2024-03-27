package main

import "mailer-server/cmd/api"

func main() {
	api.NewApiServer(":8080").Run()
}
