package producer_test

import (
	"auth_server/cmd/api/broker/channel"
	"auth_server/cmd/api/broker/producer"
	"auth_server/cmd/api/utils/connection"
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
const channel_name = "test_channel"

var prod *producer.RMQProducer

var testCreds *brokerCreds

func TestMain(m *testing.M) {
	creds, cancel := setupRabbitContainer()
	testCreds = new(brokerCreds)
	testCreds = creds
	newProd, err := producer.NewProducer().
		Host(creds.host).
		Port(creds.port).
		User(creds.user).
		Password(creds.password).
		VirtualHost(creds.vhost).
		WithChannel(
			channel_name,
			channel.NewPipeBuilder().
				Kind("direct").
				Exchange(exchange_name).
				Queue(queue_name).
				RoutingKey(routing_key),
		).Build()
	if err != nil {
		log.Fatal(err)
	}
	prod = new(producer.RMQProducer)
	prod = newProd

	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s:%s/%s",
		creds.user,
		creds.password,
		creds.host,
		creds.port,
		creds.vhost,
	)
	consumerConn, err := newConnection(dsn)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	consumerChannel, err := newConsumerChannel(consumerConn)
	if err != nil {
		log.Fatal(err)
	}

	prod.ComputeBackgroundServices()

	code := m.Run()

	time.Sleep(1 * time.Second)

	consumerChannel.Cancel(consumer_tag, false)
	consumerChannel.Close()
	consumerConn.Close()

	prod.CloseProducer()
	cancel()

	os.Exit(code)
}

func TestMakeNewProducer(t *testing.T) {
	newProd, err := producer.NewProducer().
		User(testCreds.user).
		Password(testCreds.password).
		Host(testCreds.host).
		Port(testCreds.port).
		VirtualHost(testCreds.vhost).
		WithCredentialChannel(make(chan connection.Credentials)).
		WithChannel(
			"other_channel",
			channel.NewPipeBuilder().
				Exchange("other_exchange").
				Kind("direct").
				Queue("other_queue").
				RoutingKey("other_routing_key"),
		).Build()
	if err != nil {
		t.Fatal(err)
	}

	if err := newProd.CloseProducer(); err != nil {
		t.Fatal(err)
	}
}

func TestPublishMessage(t *testing.T) {
	if err := prod.SendMessage(channel_name, []byte("halli hallo")); err != nil {
		t.Fatal(err)
	}
}

func TestPublishMessageFailed(t *testing.T) {
	err := prod.SendMessage("invalid_channel", []byte("boo wendy testaburger boo"))
	if err == nil {
		t.Fatal("should fail but got succeed")
	}

	t.Logf("got err: %s", err.Error())
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
