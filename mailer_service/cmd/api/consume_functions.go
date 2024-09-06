package api

import (
	"encoding/json"
	"fmt"
	"log"
	"mailer_service/cmd/email"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (application *app) consume(messageBus <-chan amqp.Delivery) {
	for msg := range messageBus {
		application.prepareAndSend(msg)
	}
}

func (application *app) prepareAndSend(message amqp.Delivery) {
	payload := &email.ValidationMessagePayload{}
	err := json.Unmarshal(message.Body, payload)
	if err != nil {
		log.Println("could not wrap data")
		return
	}

	kind, err := getKind(payload.Kind)
	if err != nil {
		log.Println(err.Error())
		message.Nack(false, false)
		return
	}

	message.Ack(false)

	request := email.MessageRequest{
		CorrelationID: message.CorrelationId,
		MSG: &email.Message{
			To:      payload.Email,
			Subject: kind,
			Payload: &email.ValidationMessage{
				Kind:             kind,
				Name:             payload.Name,
				ValidationNumber: payload.ValidationNumber,
			},
		},
	}

	application.mailHost.Wait.Add(1)
	application.mailerChannel <- request
}

func getKind(kind string) (string, error) {
	switch kind {
	case "verify":
		return "verify your password", nil

	case "forgotten_password":
		return "change your password", nil

	default:
		return "", fmt.Errorf("invalid email type")
	}
}
