package producer

import (
	"auth_server/cmd/api/broker/channel"
	"crypto/tls"
	"fmt"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (builder *RMQProducerBuilder) newConnection() (*amqp.Connection, error) {
	var connection *amqp.Connection
	count := 0
	var backoff time.Duration

	tlsConfig, err := builder.getTLSConfig()
	if err != nil {
		return nil, err
	}

	dsn := builder.getDSN(tlsConfig)

	for {
		conn, err := builder.createNewConnection(dsn, tlsConfig)
		if err == nil {
			connection = conn
			break
		} else {
			count++
		}

		if count > 5 {
			return nil, err
		}

		backoff = time.Duration(math.Pow(float64(count), 2)) * time.Second
		time.Sleep(backoff)
	}

	return connection, nil
}

func (builder *RMQProducerBuilder) getTLSConfig() (*tls.Config, error) {
	if builder.tlsBuilder != nil {
		return builder.tlsBuilder.Build()
	}

	return nil, nil
}

func (builder *RMQProducerBuilder) getDSN(tlsConfig *tls.Config) string {
	if tlsConfig == nil {
		return fmt.Sprintf(
			"amqp://%s:%s@%s:%s/%s",
			builder.user,
			builder.password,
			builder.host,
			builder.port,
			builder.virtualHost)
	}

	return fmt.Sprintf(
		"amqps://%s:%s@%s:%s/%s",
		builder.user,
		builder.password,
		builder.host,
		builder.port,
		builder.virtualHost)
}

func (builder *RMQProducerBuilder) createNewConnection(dsn string, tlsConfig *tls.Config) (*amqp.Connection, error) {
	if tlsConfig != nil {
		return amqp.DialTLS(dsn, tlsConfig)
	}

	return amqp.Dial(dsn)
}

func newChannels(conn *amqp.Connection, channelBuilders map[string]*channel.PipeBuilder) (map[string]*channel.Pipe, error) {
	channels := make(map[string]*channel.Pipe)
	for key, channelBuilder := range channelBuilders {
		channel, err := channelBuilder.Build(conn)
		if err != nil {
			return nil, err
		}

		channels[key] = channel
	}

	return channels, nil
}

func setLoggingChannels(builder *RMQProducerBuilder) {
	builder.WithChannel(
		"logging",
		channel.NewPipeBuilder().
			Exchange("logger_service").
			Kind("direct").
			Queue("logs").
			RoutingKey("logstore"),
	)
}
