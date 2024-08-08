package server

import (
	"listener/cmd/grpcclient"
	"listener/cmd/rabbitmq"

	"gorm.io/gorm"
)

type application struct {
	rabbitServer  *rabbitmq.RabbitServer
	grpcStructure *grpcclient.GRPCClientStructure
	database      *gorm.DB
}

func StartAndListen() error {
	app, err := setup()
	if err != nil {
		return err
	}
	defer app.grpcStructure.CloseLoggerConnection()
	defer app.grpcStructure.CloseValidationEmailConnection()

	return app.rabbitServer.Consume(app.grpcStructure, app.database)
}
