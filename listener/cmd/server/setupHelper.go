package server

import (
	"fmt"
	"listener/cmd/database"
	"listener/cmd/grpcclient"
	"listener/cmd/rabbitmq"
)

func setup() (*application, error) {
	app := &application{}

	if err := app.startRabbitServer(); err != nil {
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

	db, err := database.InitialMigration(&database.Message{})
	if err != nil {
		return nil, err
	}
	app.database = db

	grpcStructure, err := grpcclient.SetupClientStructure()
	if err != nil {
		return nil, err
	}
	app.grpcStructure = grpcStructure

	return app, nil
}

func (app *application) startRabbitServer() error {
	server, err := rabbitmq.NewRabbitServer(
		"IAmTheUser",
		"ThisIsMyPassword",
		"localhost:5672",
		"kbach",
	)
	if err != nil {
		return err
	}

	app.rabbitServer = server

	return nil
}

func (app *application) createRabbitExchange() (string, error) {
	exchange, err := app.rabbitServer.DeclareExchange(
		"users",
		"direct",
		true)
	if err != nil {
		return "", err
	}

	return exchange, err
}

func (app *application) createRabbitQueues() ([]string, error) {
	queues, err := app.rabbitServer.CreateQueues(
		true,
		false,
		"validation",
		"user",
	)
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
		return app.rabbitServer.BindQueues(
			queue,
			exchange,
			"user_validation",
			"forgotten_password",
		)

	case 1:
		return app.rabbitServer.BindQueues(
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
