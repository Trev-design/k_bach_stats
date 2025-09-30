package producer

import (
	"auth_server/cmd/api/broker/channel"
	"auth_server/cmd/api/sidecar"
	"auth_server/cmd/api/tlsconf"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RMQProducerService struct {
	RMQProducer
}

type RMQProducer struct {
	connection *amqp.Connection
	channels   map[string]*channel.Pipe
}

type RMQProducerBuilder struct {
	host        string
	port        string
	user        string
	password    string
	virtualHost string
	tlsBuilder  *tlsconf.TLSBuilder
	channels    map[string]*channel.PipeBuilder
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

func (builder *RMQProducerBuilder) Build() (*RMQProducerService, error) {
	conn, err := builder.newConnection()
	if err != nil {
		return nil, err
	}

	setLoggingChannels(builder)

	channels, err := newChannels(conn, builder.channels)
	if err != nil {
		return nil, err
	}

	return &RMQProducerService{
		RMQProducer: RMQProducer{
			connection: conn,
			channels:   channels,
		},
	}, nil
}

// registers and execute backround process/es from producer service.
func (service *RMQProducerService) HandleBackgroundProcess(req sidecar.Request) {
	defer req.Done()
	payload := req.GetPayload()
	service.SendMessage("logging", payload.GetPayloadBytes())
}

// sends a message in the background.
func (rmq *RMQProducer) SendMessage(channelName string, message []byte) error {
	channel, ok := rmq.channels[channelName]
	if !ok {
		return errors.New("something went wrong")
	}

	channel.ApplyMessage(message)

	return nil
}

// background computation of messages.
// evrey channel has its own goroutine.
func (rmq *RMQProducer) ComputeBackgroundServices() {
	for _, channel := range rmq.channels {
		go channel.ListenForMessages()
	}
}

func (rmq *RMQProducer) CloseProducer() error {
	for _, channel := range rmq.channels {
		if err := channel.CloseChannel(); err != nil {
			return err
		}
	}

	return rmq.connection.Close()
}
