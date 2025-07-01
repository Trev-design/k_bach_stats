package channel

import amqp "github.com/rabbitmq/amqp091-go"

func newChannel(conn *amqp.Connection, builder *PipeBuilder) (*amqp.Channel, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := declareExchange(channel, builder); err != nil {
		return nil, err
	}

	if err := declareQueue(channel, builder); err != nil {
		return nil, err
	}

	if err := bindQueue(channel, builder); err != nil {
		return nil, err
	}

	return channel, nil
}

func declareExchange(channel *amqp.Channel, builder *PipeBuilder) error {
	return channel.ExchangeDeclare(
		builder.exchange,
		builder.kind,
		true,
		false,
		false,
		false,
		nil,
	)
}

func declareQueue(channel *amqp.Channel, builder *PipeBuilder) error {
	if _, err := channel.QueueDeclare(
		builder.queue,
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

func bindQueue(channel *amqp.Channel, builder *PipeBuilder) error {
	return channel.QueueBind(
		builder.queue,
		builder.routingKey,
		builder.exchange,
		false,
		nil,
	)
}
