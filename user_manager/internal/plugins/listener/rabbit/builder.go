package rabbit

import (
	"errors"
	"sync"
	"user_manager/internal/core"
)

type listenerAdapterBuilder struct {
	adapter ListenerAdapter
}

func NewListenerAdapterBuilder() *listenerAdapterBuilder {
	return &listenerAdapterBuilder{
		adapter: ListenerAdapter{
			waitGroup: &sync.WaitGroup{},
		}}
}

func (builder *listenerAdapterBuilder) Connection(user, password, host, vhost string, port int) *listenerAdapterBuilder {
	connection, err := connect(user, password, host, vhost, port)
	if err != nil {
		panic("could not connect to rabbitmq")
	}

	builder.adapter.connection = connection

	return builder
}

func (builder *listenerAdapterBuilder) Channel(queue, exchange, routingKey, consumerTag string) *listenerAdapterBuilder {
	if builder.adapter.connection == nil {
		panic("connection must initialized first")
	}

	ch, err := builder.adapter.connection.Channel()
	if err != nil {
		panic("could not create channel")
	}

	if err := createQueue(ch, queue); err != nil {
		panic("could not create queue")
	}

	if err := bindQueue(
		ch,
		queue,
		exchange,
		routingKey,
	); err != nil {
		panic("could not bind queue")
	}

	builder.adapter.rabbitChannels = append(
		builder.adapter.rabbitChannels,
		&channel{
			RoutingKey:  routingKey,
			ConsumerTag: consumerTag,
			Queue:       queue,
			Chan:        ch,
		},
	)

	return builder
}

func (builder *listenerAdapterBuilder) UserStore(store core.UserManagement) *listenerAdapterBuilder {
	builder.adapter.userStore = store
	return builder
}

func (builder *listenerAdapterBuilder) SessionStore(store core.UserManagement) *listenerAdapterBuilder {
	builder.adapter.sessionStore = store
	return builder
}

func (builder *listenerAdapterBuilder) Build() (*ListenerAdapter, error) {
	if builder.adapter.connection == nil {
		return nil, errors.New("on build recognized: connection not implemented")
	}

	if len(builder.adapter.rabbitChannels) < 4 {
		return nil, errors.New("on build recognized: not enough channels")
	}

	if builder.adapter.userStore == nil {
		return nil, errors.New("on build recognized: no user store implementation")
	}

	if builder.adapter.sessionStore == nil {
		return nil, errors.New("on build recognized: no session store implemented")
	}

	return &builder.adapter, nil
}
