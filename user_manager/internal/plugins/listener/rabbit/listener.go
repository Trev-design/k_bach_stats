package rabbit

import (
	"sync"
	"user_manager/internal/core"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ListenerAdapter struct {
	connection     *amqp.Connection
	rabbitChannels []*channel
	waitGroup      *sync.WaitGroup
	userStore      core.UserManagement
	sessionStore   core.UserManagement
}

type ChannelCredentials struct {
	QueueName    string
	RoutingKey   string
	ConsumerTag  string
	ExchangeName string
}

type channel struct {
	RoutingKey  string
	ConsumerTag string
	Queue       string
	Chan        *amqp.Channel
}

type handlerFunc func(payload []byte) error

func (la *ListenerAdapter) Consume() error {
	for _, rmqChannel := range la.rabbitChannels {
		ch, err := rmqChannel.Chan.Consume(
			rmqChannel.Queue,
			rmqChannel.RoutingKey,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		go rmqChannel.consume(ch, la.userStore, la.sessionStore)
	}

	return nil
}

func (la *ListenerAdapter) CloseListener() error {
	for _, channel := range la.rabbitChannels {
		if err := channel.Chan.Close(); err != nil {
			return err
		}
	}

	if err := la.connection.Close(); err != nil {
		return err
	}

	return nil
}

func (la *ListenerAdapter) Disconnect() error {
	for _, channel := range la.rabbitChannels {
		err := channel.Chan.Cancel(channel.ConsumerTag, false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (la *ListenerAdapter) Wait() {
	la.waitGroup.Wait()
}

func (ch *channel) consume(deliveryChannel <-chan amqp.Delivery, userStore, sessionStore core.UserManagement) {
	switch ch.ConsumerTag {
	case "add_session_consumer":
		ch.consumeMessages(deliveryChannel, func(payload []byte) error {
			return sessionStore.AddUser(payload)
		})

	case "remove_session_consumer":
		ch.consumeMessages(deliveryChannel, func(payload []byte) error {
			return sessionStore.RemoveUser(payload)
		})

	case "add_user_consumer":
		ch.consumeMessages(deliveryChannel, func(payload []byte) error {
			return userStore.AddUser(payload)
		})

	case "remove_user_consumer":
		ch.consumeMessages(deliveryChannel, func(payload []byte) error {
			return userStore.RemoveUser(payload)
		})

	default:
		panic("invalid consumer tag")
	}
}

func (ch *channel) consumeMessages(deliveryChannel <-chan amqp.Delivery, handlerFn handlerFunc) {
	for message := range deliveryChannel {
		if err := handlerFn(message.Body); err != nil {
			message.Nack(false, false)
		}

		message.Ack(false)
	}
}
