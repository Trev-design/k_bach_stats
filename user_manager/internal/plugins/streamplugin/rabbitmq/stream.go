package rabbitmq

import (
	"errors"
	"user_manager/types"

	"github.com/google/uuid"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/message"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

type StreamAdapter struct {
	environment         *stream.Environment
	invitationProducer  *stream.Producer
	joinRequestProducer *stream.Producer
}

func NewStreamAdapter() (*StreamAdapter, error) {
	return makeEnv()
}

func (sa *StreamAdapter) SendStream(data *types.StreamPayload) error {
	switch data.Kind {
	case "join_request":
		return send(sa.joinRequestProducer, data.Payload)

	case "invitation":
		return send(sa.invitationProducer, data.Payload)

	default:
		return errors.New("invalid kind")
	}
}

func (sa *StreamAdapter) Close() error {
	if err := sa.environment.Close(); err != nil {
		return err
	}

	if err := sa.invitationProducer.Close(); err != nil {
		return err
	}

	return sa.joinRequestProducer.Close()
}

func send(producer *stream.Producer, data [][]byte) error {
	batchData := make([]message.StreamMessage, len(data))

	for index, payload := range data {
		message := amqp.NewMessage(payload)
		props := &amqp.MessageProperties{CorrelationID: uuid.New().String()}
		message.Properties = props
		message.SetPublishingId(int64(index))

		batchData[index] = message
	}

	return sendBatches(producer, batchData)
}

func sendBatches(producer *stream.Producer, data []message.StreamMessage) error {
	batchSize := 20

	for index := 0; index+batchSize < len(data); index += batchSize {
		end := index + batchSize

		if end > len(data) {
			end = len(data)
		}

		batch := data[index:end]
		if err := producer.BatchSend(batch); err != nil {
			return err
		}
	}

	return nil
}
