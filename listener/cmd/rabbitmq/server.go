package rabbitmq

import (
	"encoding/json"
	"listener/cmd/grpcclient"
	"listener/cmd/messagetypes"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type consumerFn func(
	message *messagetypes.Message,
	clients *grpcclient.GRPCClientStructure,
	db *gorm.DB,
) error

func (server *RabbitServer) Consume(clients *grpcclient.GRPCClientStructure, db *gorm.DB) error {
	validationBus, err := server.createConsumer(
		"validation",
		"validation_mail_service",
	)
	if err != nil {
		return err
	}

	modifyUserBus, err := server.createConsumer(
		"user",
		"user_modify_service",
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)

	go server.consumeQueue(validationBus, computeEmailMessage, clients, db)
	go server.consumeQueue(modifyUserBus, computeModifyUserMessage, clients, db)

	<-forever

	log.Println("Ciao Miao Miao")

	return nil
}

func (server *RabbitServer) createConsumer(queue, consumerTag string) (<-chan amqp.Delivery, error) {
	return server.channel.Consume(
		queue,
		consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
}

func (server *RabbitServer) consumeQueue(
	messages <-chan amqp.Delivery,
	fun consumerFn,
	clients *grpcclient.GRPCClientStructure,
	db *gorm.DB,
) {

	for message := range messages {
		body := new(messagetypes.Message)
		err := json.Unmarshal(message.Body, body)
		if err != nil {
			log.Printf("could not parse json: %v\n", err)
		}

		if err := fun(body, clients, db); err != nil {
			log.Printf("could not send data: %v", err)
		} else {
			message.Ack(false)
		}
	}
}
