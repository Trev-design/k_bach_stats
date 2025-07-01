package rabbit

import (
	"encoding/json"
	"errors"
	"mail_server/cmd/api/mailer"

	amqp "github.com/rabbitmq/amqp091-go"
)

func consume(
	consumeChannel <-chan amqp.Delivery,
	messageChannel chan<- mailer.MessageRequest,
) {
	for msg := range consumeChannel {
		creds := new(mailer.MessagePayload)
		if err := json.Unmarshal(msg.Body, creds); err != nil {
			msg.Nack(false, false)
		}

		kind, err := getKind(creds.Kind)
		if err != nil {
			msg.Nack(false, false)
		}

		request := mailer.MessageRequest{
			CorrelationID: msg.CorrelationId,
			Message: &mailer.Message{
				To:      creds.Email,
				Subject: kind,
				Payload: &mailer.ValidationMessage{
					ValidationNumber: creds.ValidationNumber,
					Name:             creds.Name,
					Kind:             creds.Email,
				},
			},
		}

		messageChannel <- request
		msg.Ack(false)
	}
}

func getKind(kind string) (string, error) {
	switch kind {
	case "verify":
		return "verify your password", nil

	case "change_password":
		return "change your password", nil

	default:
		return "", errors.New("invalid message type")
	}
}
