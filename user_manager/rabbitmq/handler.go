package rabbitmq

import (
	"errors"
	"log"
	"sync"
	"user_manager/database"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitConsumerStructure struct {
	Conn             *amqp.Connection
	StartSessionChan *ChannelInstance
	StopSessionChan  *ChannelInstance
	AddUserChan      *ChannelInstance
	RemoveUserChan   *ChannelInstance
	Wait             *sync.WaitGroup
}

type ChannelInstance struct {
	Channel     *amqp.Channel
	QueueName   string
	ConsumerTag string
}

type RabbitConsumer interface {
	ComputeMessages(job string, userHandler database.UserHandler, errorChan chan error)
	Close() error
}

type userHandlerFn func([]byte) error

func (structure *RabbitConsumerStructure) ComputeMessages(job string, userHandler database.UserHandler, errorChan chan error) {
	switch job {
	case "start_session":
		structure.StartSessionChan.consume(
			errorChan,
			structure.Wait,
			func(payload []byte) error { return userHandler.AddUser(payload) },
		)

	case "stop_session":
		structure.StopSessionChan.consume(
			errorChan,
			structure.Wait,
			func(payload []byte) error { return userHandler.RemoveUser(payload) },
		)

	case "add_user":
		structure.AddUserChan.consume(
			errorChan,
			structure.Wait,
			func(payload []byte) error { return userHandler.AddUser(payload) },
		)
	case "remove_user":
		structure.RemoveUserChan.consume(
			errorChan,
			structure.Wait,
			func(payload []byte) error { return userHandler.RemoveUser(payload) },
		)
	default:
		errorChan <- errors.New("invalid job request")
	}
}

func (instance *ChannelInstance) consume(errorChannel chan error, wait *sync.WaitGroup, fn userHandlerFn) {
	channel, err := instance.Channel.Consume(
		instance.QueueName,
		instance.ConsumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		errorChannel <- err
	}

	log.Println("ready for consuming messages")

	for message := range channel {
		wait.Add(1)
		go func() {
			if err := instance.computeMessage(wait, fn, message.Body); err != nil {
				instance.Channel.Nack(message.DeliveryTag, false, false)
				log.Printf("the err: %v", err)
				errorChannel <- err
			}

			instance.Channel.Ack(message.DeliveryTag, false)
		}()
	}
}

func (structure *RabbitConsumerStructure) Close() error {
	structure.AddUserChan.Channel.Cancel(structure.AddUserChan.ConsumerTag, false)
	structure.RemoveUserChan.Channel.Cancel(structure.RemoveUserChan.ConsumerTag, false)
	structure.StartSessionChan.Channel.Cancel(structure.StartSessionChan.ConsumerTag, false)
	structure.StopSessionChan.Channel.Cancel(structure.StopSessionChan.ConsumerTag, false)

	structure.Wait.Wait()
	structure.AddUserChan.Channel.Close()
	structure.RemoveUserChan.Channel.Close()
	structure.StartSessionChan.Channel.Close()
	structure.StartSessionChan.Channel.Close()
	structure.Conn.Close()

	return nil
}

func (instance *ChannelInstance) computeMessage(wait *sync.WaitGroup, fn userHandlerFn, payload []byte) error {
	defer wait.Done()
	log.Printf("the very cool payload of the message is %v", payload)
	return fn(payload)
}
