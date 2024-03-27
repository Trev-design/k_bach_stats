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
	go server.listenForShutDown()

	srv := &http.Server{
		Addr:    server.port,
		Handler: server.routes(),
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
