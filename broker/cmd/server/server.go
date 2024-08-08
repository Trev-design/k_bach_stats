package server

import (
	"broker-server/cmd/rabbit"

	"github.com/gofiber/fiber/v2"
)

type application struct {
	httpServer   *fiber.App
	rabbitClient *rabbit.RabbitClient
}

func StartAndListen() error {

	app, err := setup()
	if err != nil {
		return err
	}

	defer app.rabbitClient.CloseConnection()
	defer app.rabbitClient.CloseChannel()

	app.routes()

	return app.httpServer.Listen(":4001")
}
