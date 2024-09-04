package api

import (
	"fmt"
	"log"
	"mailer_service/cmd/email"
	"mailer_service/cmd/rabbitmq"
	"sync"
)

type app struct {
	rmqClient     *rabbitmq.RabbitClient
	mailHost      *email.MailHost
	wait          *sync.WaitGroup
	mailerChannel chan email.MessageRequest
}

func StartAndListen() error {
	application, err := newApp()
	if err != nil {
		return fmt.Errorf("could not start server. error: %v", err)
	}
	application.wait.Add(1)

	consumerChannel, err := application.rmqClient.CreateConsumer("verification_email")
	if err != nil {
		return fmt.Errorf("could not start consumer. error: %v", err)
	}

	log.Println("server started")

	go application.consume(consumerChannel)
	go application.ListenForEmails()
	go application.listenForShutDown()

	application.wait.Wait()

	return nil
}
