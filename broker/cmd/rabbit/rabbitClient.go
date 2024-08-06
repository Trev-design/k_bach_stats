package rabbit

import (
	messagetypes "broker-server/cmd/message_types"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func (client *RabbitClient) PublishMessage(ctx context.Context, exchange, key, payloadHeader string, payload any) error {
	msg := &messagetypes.Message{
		Type:    payloadHeader,
		Payload: payload,
	}

	message, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return client.channel.PublishWithContext(
		ctx,
		exchange,
		key,
		true,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         message,
		},
	)
}

func NewrabbitClient(username, password, host, vhost string) (*RabbitClient, error) {
	connection, err := connectToRabbit(username, password, host, vhost)
	if err != nil {
		return nil, fmt.Errorf("failed to start rabbit connection because of this error: %v", err)
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to start rabbitmq channel because of tis following error %v", err)
	}

	return &RabbitClient{
		connection: connection,
		channel:    channel,
	}, nil
}

func (client *RabbitClient) CreateQueues(durable, autoDelete bool, queueNames ...string) ([]string, error) {
	var names []string

	for _, name := range queueNames {
		if _, err := client.channel.QueueDeclare(
			name,
			durable,
			autoDelete,
			false,
			false,
			nil,
		); err != nil {
			return nil, fmt.Errorf(
				"failed to create queue with name of %s with error: %v",
				name,
				err,
			)
		}

		names = append(names, name)
	}

	return names, nil
}

func (client *RabbitClient) DeclareExchange(name, kind string, durable bool) (string, error) {
	if err := client.channel.ExchangeDeclare(
		name,
		kind,
		durable,
		false,
		false,
		false,
		nil,
	); err != nil {
		return "", err
	}

	return name, nil
}

func (client *RabbitClient) BindQueues(name, exchange string, keys ...string) error {
	for _, key := range keys {
		if err := client.channel.QueueBind(
			name,
			key,
			exchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf(
				"could not create binding of queue %s with key %s with this error: %v",
				name,
				key,
				err,
			)
		}
	}
	return nil
}

func (client *RabbitClient) CloseConnection() error {
	return client.connection.Close()
}

func (client *RabbitClient) CloseChannel() error {
	return client.channel.Close()
}

func connectToRabbit(username, password, host, vhost string) (*amqp.Connection, error) {
	var counts int
	var backOff = 1 * time.Second
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s/%s", username, password, host, vhost))
		if err != nil {
			counts++
		} else {
			connection = conn
			break
		}

		if counts > 5 {
			return nil, err
		}

		backOff = time.Duration(math.Pow(float64(counts), 2)) * time.Second
		time.Sleep(backOff)
	}

	return connection, nil
}
