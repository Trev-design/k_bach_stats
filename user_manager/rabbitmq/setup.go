package rabbitmq

import (
	"fmt"
	"math"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Setup() (*RabbitConsumerStructure, error) {
	conn, err := connect("IAmTheUser", "ThisIsMyPassword", "localhost", "kbach", 5672)
	if err != nil {
		return nil, err
	}

	addSession, err := setupChannels("add_session_consumer", "start_user_session", "session", "send_session_credentials", conn)
	if err != nil {
		return nil, err
	}

	removeSession, err := setupChannels("remove_session_consumer", "stop_user_session", "session", "remove_user_session", conn)
	if err != nil {
		return nil, err
	}

	addUser, err := setupChannels("add_user_consumer", "add_account", "account", "add_account_request", conn)
	if err != nil {
		return nil, err
	}

	removeUser, err := setupChannels("remove_user_consumer", "delete_account", "account", "delete_account_request", conn)
	if err != nil {
		return nil, err
	}

	return &RabbitConsumerStructure{
		Conn:             conn,
		StartSessionChan: addSession,
		StopSessionChan:  removeSession,
		AddUserChan:      addUser,
		RemoveUserChan:   removeUser,
		Wait:             &sync.WaitGroup{},
	}, nil
}

func setupChannels(consumertag, queuename, exchangename, routingkey string, conn *amqp.Connection) (*ChannelInstance, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := CreateQueue(channel, queuename, true, false); err != nil {
		return nil, err
	}

	if err := BindQueue(queuename, exchangename, routingkey, channel); err != nil {
		return nil, err
	}

	return &ChannelInstance{
		Channel:     channel,
		ConsumerTag: consumertag,
	}, nil
}

func CreateQueue(channel *amqp.Channel, name string, durable, autDelete bool) error {
	if _, err := channel.QueueDeclare(
		name,
		durable,
		autDelete,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func BindQueue(queuename, exchangename, routingkey string, channel *amqp.Channel) error {
	return channel.QueueBind(
		queuename,
		routingkey,
		exchangename,
		false,
		nil,
	)
}

func connect(user, password, host, vhost string, port int) (*amqp.Connection, error) {
	count := 0
	backoff := 1 * time.Second
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

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		time.Sleep(backoff)
	}

	return connection, nil
}
