package server

import (
	"listener/cmd/rabbitmq"

	"gorm.io/gorm"
)

type application struct {
	rabbitServer *rabbitmq.RabbitServer
	database     *gorm.DB
}

func StartAndListen() error {
	app, err := setup()
	if err != nil {
		return err
	}

	return app.rabbitServer.Consume()
}
