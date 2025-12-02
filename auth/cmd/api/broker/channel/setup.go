package channel

import (
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Pipe struct {
	channel     *amqp.Channel
	messagePipe chan []byte
	routingKey  string
	exchange    string
	waitgroup   *sync.WaitGroup
}

type PipeBuilder struct {
	exchange   string
	queue      string
	routingKey string
	kind       string
}

// We use a builder to initialize the channel structure of the rabbitmq service
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
		waitgroup:   &sync.WaitGroup{},
		messagePipe: make(chan []byte, 100),
		routingKey:  builder.routingKey,
		exchange:    builder.exchange,
	}, nil
}

// Puts a message in the messagechannel we use a []byte channel because of simplicity and thread safetyness
func (pipe *Pipe) ApplyMessage(message []byte) {
	pipe.waitgroup.Add(1)
	pipe.messagePipe <- message
}

// Computes the messages comming out of the messagechannel.
// This function needs to start in a goroutine to run in the background
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

		pipe.waitgroup.Done()
	}
}

// Cleans up the channel.
// This function should only be called by server shutdown or some restart mechanisms
func (pipe *Pipe) CloseChannel() error {
	close(pipe.messagePipe)
	pipe.waitgroup.Wait()
	return pipe.channel.Close()
}
