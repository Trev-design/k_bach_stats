package channel

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Pipe struct {
	channel     *amqp.Channel
	messagePipe chan []byte
	routingKey  string
	exchange    string
}

type PipeBuilder struct {
	exchange   string
	queue      string
	routingKey string
	kind       string
}

func NewPipeBuilder() *PipeBuilder {
	return &PipeBuilder{}
}

func (builder *PipeBuilder) Exchange(exchange string) *PipeBuilder {
	builder.exchange = exchange
	return builder
}

func (builder *PipeBuilder) Queue(queue string) *PipeBuilder {
	builder.queue = queue
	return builder
}

func (builder *PipeBuilder) Kind(kind string) *PipeBuilder {
	builder.kind = kind
	return builder
}

func (builder *PipeBuilder) RoutingKey(routingKey string) *PipeBuilder {
	builder.routingKey = routingKey
	return builder
}

func (builder *PipeBuilder) Build(conn *amqp.Connection) (*Pipe, error) {
	channel, err := newChannel(conn, builder)
	if err != nil {
		return nil, err
	}

	return &Pipe{
		channel:     channel,
		messagePipe: make(chan []byte, 100),
		routingKey:  builder.routingKey,
		exchange:    builder.exchange,
	}, nil
}

func (pipe *Pipe) ApplyMessage(message []byte) {
	pipe.messagePipe <- message
}

func (pipe *Pipe) ListenForMessages() {
	for message := range pipe.messagePipe {
		if err := pipe.channel.Publish(
			pipe.exchange,
			pipe.routingKey,
			true,
			false,
			amqp.Publishing{
				ContentType:  "text/plain",
				DeliveryMode: amqp.Persistent,
				Body:         message,
			},
		); err != nil {
			log.Println(err)
		}
	}
}

func (pipe *Pipe) CloseChannel() error {
	close(pipe.messagePipe)
	return pipe.channel.Close()
}
