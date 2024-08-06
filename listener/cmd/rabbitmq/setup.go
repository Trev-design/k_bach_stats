package rabbitmq

import (
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitServer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
}

func NewRabbitServer(username, password, host, vhost string) (*RabbitServer, error) {
	connection, err := connectToRabbit(username, password, host, vhost)
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitServer{
		connection: connection,
		channel:    channel,
	}, nil
}

func connectToRabbit(username, password, host, vhost string) (*amqp.Connection, error) {
	var counts int
	var backoff = 1 * time.Second
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

		backoff = time.Duration(math.Pow(float64(counts), 2))
		time.Sleep(backoff)
	}

	return connection, nil
}

func (server *RabbitServer) CreateQueues(durable, autoDelete bool, queueNames ...string) ([]string, error) {
	var names []string
	for _, name := range queueNames {
		if _, err := server.channel.QueueDeclare(
			name,
			durable,
			autoDelete,
			false,
			false,
			nil,
		); err != nil {
			return nil, fmt.Errorf("could not create queue %s: %v", name, err)
		}

		names = append(names, name)
	}

	return names, nil
}

func (server *RabbitServer) DeclareExchange(name, kind string, durable bool) (string, error) {
	if err := server.channel.ExchangeDeclare(
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

func (server *RabbitServer) BindQueues(name, exchange string, keys ...string) error {
	for _, key := range keys {
		if err := server.channel.QueueBind(
			name,
			key,
			exchange,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("could not bind queue %s on key %s: %v", name, key, err)
		}
	}

	return nil
}

func (server *RabbitServer) CloseConnection() error {
	return server.connection.Close()
}

func (server *RabbitServer) CloseChannel() error {
	return server.channel.Close()
}
