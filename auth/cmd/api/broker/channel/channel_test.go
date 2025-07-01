package channel_test

import (
	"auth_server/cmd/api/broker/channel"
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type brokerCreds struct {
	user     string
	password string
	host     string
	port     string
	vhost    string
}

const routing_key = "test_key"
const exchange_name = "test_exchange"
const queue_name = "test_queue"
const consumer_tag = "test_tag"

var testChannel *channel.Pipe
var testConnection *amqp.Connection

func TestMain(m *testing.M) {
	creds, cancel := setupRabbitContainer()

	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", creds.user, creds.password, creds.host, creds.port, creds.vhost)
	conn, err := newConnection(dsn)
	if err != nil {
		log.Fatal(err)
	}
	testConnection = conn
	consumerConn, err := newConnection(dsn)
	if err != nil {
		log.Fatal(err)
	}

	channel, err := channel.NewPipeBuilder().
		Exchange(exchange_name).
		RoutingKey(routing_key).
		Queue(queue_name).
		Kind("direct").
		Build(conn)
	if err != nil {
		log.Fatal(err)
	}

	testChannel = channel
	go testChannel.ListenForMessages()

	consumeChannel, err := newConsumerChannel(consumerConn)
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	code := m.Run()
	time.Sleep(1 * time.Second)
	channel.CloseChannel()
	consumeChannel.Cancel(consumer_tag, false)
	conn.Close()
	consumerConn.Close()
	cancel()
	os.Exit(code)
}

func TestMakeNewChannel(t *testing.T) {
	newChannel, err := channel.NewPipeBuilder().
		Exchange(exchange_name).
		Kind("direct").
		Queue("custom_queue").
		RoutingKey("custom_routing_key").
		Build(testConnection)
	if err != nil {
		t.Fatal(err)
	}

	newChannel.CloseChannel()
}

func TestSendMessage(t *testing.T) {
	testChannel.ApplyMessage([]byte("halli hallo"))
}

func setupRabbitContainer() (*brokerCreds, func()) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "rabbitmq:4-management",
		ExposedPorts: []string{"5672/tcp", "15672/tcp"},
		Env: map[string]string{
			"RABBITMQ_DEFAULT_USER":  "test_user",
			"RABBITMQ_DEFAULT_PASS":  "testpass",
			"RABBITMQ_DEFAULT_VHOST": "test_vhost",
		},
		WaitingFor: wait.ForLog("Server startup complete"),
		Mounts:     testcontainers.ContainerMounts{},
	}
	instance, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req, Started: true,
		},
	)
	if err != nil {
		log.Fatal("failed to start container")
	}

	host, err := instance.Host(ctx)
	if err != nil {
		log.Fatal("failed to fetch host")
	}

	port, err := instance.MappedPort(ctx, "5672/tcp")
	if err != nil {
		log.Fatal("failed to fetche mapped port")
	}

	time.Sleep(2 * time.Second)

	return &brokerCreds{
		user:     "test_user",
		password: "testpass",
		host:     host,
		port:     port.Port(),
		vhost:    "test_vhost",
	}, func() { instance.Terminate(ctx) }
}

func newConnection(dsn string) (*amqp.Connection, error) {
	return amqp.Dial(dsn)
}

func newConsumerChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	consumerChannel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	messages, err := consumerChannel.Consume(queue_name, consumer_tag, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go consumeMessages(messages)

	return consumerChannel, nil
}

func consumeMessages(messages <-chan amqp.Delivery) {
	for message := range messages {
		log.Println(string(message.Body))
		message.Ack(false)
	}
}
