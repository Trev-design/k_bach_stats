package api

import (
	"mailer_service/cmd/email"
	"mailer_service/cmd/rabbitmq"
	"sync"
)

func newApp() (*app, error) {
	rabbitClient, err := rabbitmq.NewClient(
		"IAmTheUser",
		"ThisIsMyPassword",
		"localhost",
		"kbach",
		5672,
	)
	if err != nil {
		return nil, err
	}

	if err := initVerifyEmailQueue(rabbitClient); err != nil {
		return nil, err
	}

	return &app{
		rmqClient:     rabbitClient,
		mailHost:      email.NewValidationMailer(),
		mailerChannel: make(chan email.MessageRequest, 100),
		wait:          &sync.WaitGroup{},
	}, nil
}

func initVerifyEmailQueue(client *rabbitmq.RabbitClient) error {
	if err := client.DeclareExchange(
		"verify",
		"direct",
		true,
	); err != nil {
		return err
	}

	if err := client.CreateQueue(
		"verification_email",
		true,
		false,
	); err != nil {
		return err
	}

	keys := []string{"send_verify_email", "send_forgotten_password_email"}
	for _, key := range keys {
		if err := client.BindQueue(
			"verification_email",
			key,
			"verify",
		); err != nil {
			return err
		}
	}

	return nil
}
