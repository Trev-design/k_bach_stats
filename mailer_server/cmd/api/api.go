package api

import (
	"log"
	"mailer-server/cmd/email"
	"net/http"
	"sync"
)

type app struct {
	port   string
	mailer *email.Mail
	wait   *sync.WaitGroup
}

func NewApiServer(listen string) *app {
	return &app{port: listen, wait: &sync.WaitGroup{}}
}

func (server *app) Run() {
	server.mailer = email.CreateMail(server.wait)

	go server.mailer.ListenForMail()

	err := http.ListenAndServe(server.port, server.routes())
	if err != nil {
		log.Fatal(err)
	}
}
