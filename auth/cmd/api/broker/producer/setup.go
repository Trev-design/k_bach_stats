package producer

import (
	"auth_server/cmd/api/broker/channel"
	"auth_server/cmd/api/tlsconf"
	"auth_server/cmd/api/utils/connection"
	"errors"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQProducer struct {
	builder            *RMQProducerBuilder
	credentialsChannel chan connection.Credentials
	conn               *connection.Handler[RMQConnection]
}

type RMQConnection struct {
	connection *amqp.Connection
	channels   map[string]*channel.Pipe
	waitgroup  *sync.WaitGroup
}

type RMQProducerBuilder struct {
	host        string
	port        string
	user        string
	password    string
	virtualHost string
	tlsBuilder  *tlsconf.TLSBuilder
	channels    map[string]*channel.PipeBuilder
	pipe        chan connection.Credentials
}

func NewProducer() *RMQProducerBuilder {
	return &RMQProducerBuilder{
		channels: make(map[string]*channel.PipeBuilder),
	}
}

func (builder *RMQProducerBuilder) Host(host string) *RMQProducerBuilder {
	builder.host = host
	return builder
}

func (builder *RMQProducerBuilder) Port(port string) *RMQProducerBuilder {
	builder.port = port
	return builder
}

func (builder *RMQProducerBuilder) User(user string) *RMQProducerBuilder {
	builder.user = user
	return builder
}

func (builder *RMQProducerBuilder) Password(password string) *RMQProducerBuilder {
	builder.password = password
	return builder
}

func (builder *RMQProducerBuilder) VirtualHost(virtualHost string) *RMQProducerBuilder {
	builder.virtualHost = virtualHost
	return builder
}

// for optional tls support we have a tls builder where you can parse your tls configurations
func (builder *RMQProducerBuilder) WithTLS(tlsBuilder *tlsconf.TLSBuilder) *RMQProducerBuilder {
	builder.tlsBuilder = tlsBuilder
	return builder
}

// here you can setup your channels.
// because there can more then one channels, this method can be used more than one time
func (builder *RMQProducerBuilder) WithChannel(
	channelName string,
	channelBuilder *channel.PipeBuilder,
) *RMQProducerBuilder {
	if _, ok := builder.channels[channelName]; channelName == "logging" && ok {
		return builder
	}

	builder.channels[channelName] = channelBuilder
	return builder
}

func (builder *RMQProducerBuilder) WithCredentialChannel(pipe chan connection.Credentials) *RMQProducerBuilder {
	builder.pipe = pipe
	return builder
}

func (builder *RMQProducerBuilder) BuildConnection() (connection.Connection[RMQConnection], error) {
	conn, err := builder.newConnection()
	if err != nil {
		return nil, err
	}

	channels, err := newChannels(conn, builder.channels)
	if err != nil {
		return nil, err
	}

	return &RMQConnection{
		connection: conn,
		channels:   channels,
		waitgroup:  &sync.WaitGroup{},
	}, nil
}

func (builder *RMQProducerBuilder) Build() (*RMQProducer, error) {
	conn, err := connection.NewBuilder(builder).Build()
	if err != nil {
		return nil, err
	}

	return &RMQProducer{
		conn:               conn,
		credentialsChannel: builder.pipe,
		builder:            builder,
	}, nil
}

// sends a message in the background.
func (rmq *RMQProducer) SendMessage(channelName string, message []byte) error {
	conn := rmq.conn.Get()
	conn.waitgroup.Add(1)
	defer conn.waitgroup.Done()
	channel, ok := conn.channels[channelName]
	if !ok {
		return errors.New("something went wrong")
	}

	channel.ApplyMessage(message)

	return nil
}

// background computation of messages.
// evrey channel has its own goroutine.
func (rmq *RMQProducer) ComputeBackgroundServices() {
	conn := rmq.conn.Get()
	for _, channel := range conn.channels {
		go channel.ListenForMessages()
	}
}

func (rmq *RMQProducer) CloseProducer() error {
	conn := rmq.conn.Get()
	conn.Wait()
	return conn.Close()
}
