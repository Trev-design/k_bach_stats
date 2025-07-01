package rabbit

import (
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connect(user, password, host, port, vhost string) (*amqp.Connection, error) {
	count := 0
	backoff := 1 * time.Second
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", user, password, host, port, vhost)
	var connection *amqp.Connection

	for {
		conn, err := amqp.Dial(dsn)
		if err != nil {
			count++
		} else {
			connection = conn
			break
		}

		if count > 5 {
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		time.Sleep(backoff)
	}

	return connection, nil
}

func makeBindedQueue(channel *amqp.Channel, exchange, kind, queue, key string) error {
	if err := makeExchange(channel, exchange, kind); err != nil {
		return err
	}

	if err := createQueue(channel, queue); err != nil {
		return err
	}

	return bindQueue(channel, queue, exchange, key)
}

func makeExchange(channel *amqp.Channel, exchange, kind string) error {
	return channel.ExchangeDeclare(
		exchange,
		kind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func createQueue(channel *amqp.Channel, queue string) error {
	if _, err := channel.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func bindQueue(channel *amqp.Channel, queue, exchange, key string) error {
	return channel.QueueBind(
		queue,
		key,
		exchange,
		false,
		nil,
	)
}

func (srv *Rabbit) createConsumer(queue string) (<-chan amqp.Delivery, error) {
	return srv.channel.Consume(
		queue,
		srv.tag,
		false,
		false,
		false,
		false,
		nil,
	)
}
