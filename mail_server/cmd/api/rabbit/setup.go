package rabbit

import (
	"mail_server/cmd/api/mailer"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Rabbit struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	tag     string
	queue   string
	wait    *sync.WaitGroup
}

func NewClient(user, password, host, port, vhost string) (*Rabbit, error) {
	conn, err := connect(user, password, host, port, vhost)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := makeBindedQueue(
		channel,
		"mails",
		"direct",
		"verify_mails",
		"send_verify_mail",
	); err != nil {
		return nil, err
	}

	return &Rabbit{
		conn:    conn,
		channel: channel,
		queue:   "verify_mails",
		tag:     "verify_emails",
		wait:    &sync.WaitGroup{},
	}, nil
}

func (srv *Rabbit) CloseRabbit() error {
	if err := srv.channel.Cancel(srv.tag, false); err != nil {
		return err
	}
	srv.wait.Wait()
	if err := srv.channel.Close(); err != nil {
		return err
	}

	return srv.conn.Close()
}

func (srv *Rabbit) StartConsuming(messageChannel chan<- mailer.MessageRequest) {
	consumeChannel, err := srv.createConsumer(srv.queue)
	if err != nil {
		panic("unable to start consumer")
	}

	consume(consumeChannel, messageChannel)
}
