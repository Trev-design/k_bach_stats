package rabbitmq

import (
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
}

func NewClient(user, password, host, vhost string, port int) (*RabbitClient, error) {
	connection, err := connect(user, password, host, vhost, port)
	if err != nil {
		return nil, err
	}

	channel, err := connection.Channel()
	if err != nil {
		return nil, err
	}

	return &RabbitClient{
		Connection: connection,
		Channel:    channel,
	}, nil
}

func (client *RabbitClient) CloseConnection() error {
	return client.Connection.Close()
}

func (client *RabbitClient) CloseChannel() error {
	return client.Channel.Close()
}

func (client *RabbitClient) DeclareExchange(name, kind string, durable bool) error {
	return client.Channel.ExchangeDeclare(
		name,
		kind,
		durable,
		false,
		false,
		false,
		nil,
	)
}

func (client *RabbitClient) CreateQueue(name string, durable, autoDelete bool) error {
	if _, err := client.Channel.QueueDeclare(
		name,
		durable,
		autoDelete,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (client *RabbitClient) BindQueue(queueName, routinKey, exchangeName string) error {
	return client.Channel.QueueBind(
		queueName,
		routinKey,
		exchangeName,
		false,
		nil,
	)
}

func (client *RabbitClient) CreateConsumer(queueName string) (<-chan amqp.Delivery, error) {
	return client.Channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
}

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
