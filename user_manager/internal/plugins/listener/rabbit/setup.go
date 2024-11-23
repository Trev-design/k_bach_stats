package rabbit

import (
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func connect(user, password, host, vhost string, port int) (*amqp.Connection, error) {
	count := 0
	backOff := 1 * time.Second
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, password, host, port, vhost)
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

		backOff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		time.Sleep(backOff)
	}

	return connection, nil
}

func createQueue(rmqChannel *amqp.Channel, queueName string) error {
	if _, err := rmqChannel.QueueDeclare(
		queueName,
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

func bindQueue(rmqChannel *amqp.Channel, queueName, exchangeName, routingKey string) error {
	return rmqChannel.QueueBind(
		queueName,
		routingKey,
		exchangeName,
		false,
		nil,
	)
}
