package main

import (
	"mail_server/cmd/api/mailer"
	"mail_server/cmd/api/rabbit"
	"os"
	"sync"
)

type app struct {
	rmqSrv    *rabbit.Rabbit
	mailSrv   *mailer.Mailhost
	waitgroup *sync.WaitGroup
}

func main() {
	user := os.Getenv("RABBIT_USER")
	password := os.Getenv("RABBIT_PASSWORD")
	host := os.Getenv("RABBIT_HOST")
	port := os.Getenv("RABBIT_PORT")
	vhost := os.Getenv("RABBIT_V_HOST")

	app := new(app)
	app.waitgroup = &sync.WaitGroup{}
	app.waitgroup.Add(1)
	mailerChannel := make(chan mailer.MessageRequest, 100)
	mailSrv := mailer.NewMailHost(mailerChannel)
	app.mailSrv = mailSrv

	rmqSrv, err := rabbit.NewClient(user, password, host, port, vhost)
	if err != nil {
		panic("could not start rabbit client")
	}
	app.rmqSrv = rmqSrv

	go app.rmqSrv.StartConsuming(mailerChannel)
	go app.mailSrv.ListenForEmails()
	go app.listenForShutDown()

	app.waitgroup.Wait()
	close(mailerChannel)
}
