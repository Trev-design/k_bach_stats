package server

import (
	"broker-server/cmd/rabbit"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func setup() (*application, error) {
	app := &application{}
	app.httpServer = fiber.New()

	if err := app.startRabbitClient(); err != nil {
		return nil, err
	}

	exchange, err := app.createRabbitExchange()
	if err != nil {
		return nil, err
	}

	queues, err := app.createRabbitQueues()
	if err != nil {
		return nil, err
	}

	if err := app.bindAllRabbitQueues(queues, exchange); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *application) startRabbitClient() error {
	client, err := rabbit.NewrabbitClient(
		"IAmTheUser",
		"ThisIsMyPassword",
		"localhost:5672",
		"kbach",
	)
	if err != nil {
		return err
	}
	app.rabbitClient = client

	return nil
}

func (app *application) createRabbitExchange() (string, error) {
	exchange, err := app.rabbitClient.DeclareExchange(
		"users",
		"direct",
		true)
	if err != nil {
		return "", err
	}

	return exchange, err
}

func (app *application) createRabbitQueues() ([]string, error) {
	queues, err := app.rabbitClient.CreateQueues(
		true,
		false,
		"validation",
		"user",
		"mailing")
	if err != nil {
		return nil, err
	}

	return queues, nil
}

func (app *application) bindAllRabbitQueues(queues []string, exchange string) error {
	for index, queue := range queues {
		return app.bindRabbitQueue(index, queue, exchange)
	}

	return nil
}

func (app *application) bindRabbitQueue(index int, queue, exchange string) error {
	switch index {
	case 0:
		return app.rabbitClient.BindQueues(
			queue,
			exchange,
			"user_validation",
			"forgotten_password",
		)

	case 1:
		return app.rabbitClient.BindQueues(
			queue,
			exchange,
			"create_user",
			"update_user",
			"delete_user",
		)

	default:
		return fmt.Errorf("index out of bounds")
	}
}
