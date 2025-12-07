package connection

import (
	"fmt"
	"log"
	"sync/atomic"
)

type SuccessMessage struct {
	Provider string
	Status   string
	Message  string
}

type SaltMessage struct {
	Provider string
	Salt     string
}

type Credentials struct {
	UserName string `json:"user_name"`
	Password string `jsone:"password"`
}

type Connection[Conn any] interface {
	Wait()
	Close() error
	Conn() *Conn
}

type Builder[Conn any] struct {
	connectionBuilder builder[Conn]
}

type builder[Conn any] interface {
	BuildConnection() (Connection[Conn], error)
}

type slot[Conn any] struct {
	connection Connection[Conn]
}

type Handler[Conn any] struct {
	slot atomic.Pointer[slot[Conn]]
}

func NewBuilder[Conn any](builder builder[Conn]) *Builder[Conn] {
	return &Builder[Conn]{connectionBuilder: builder}
}

func (builder *Builder[Conn]) Build() (*Handler[Conn], error) {
	conn, err := builder.connectionBuilder.BuildConnection()
	if err != nil {
		return nil, err
	}

	handler := &Handler[Conn]{
		slot: atomic.Pointer[slot[Conn]]{},
	}

	handler.slot.Store(&slot[Conn]{connection: conn})

	return handler, nil
}

func (handler *Handler[Conn]) Get() *Conn {
	return handler.slot.Load().connection.Conn()
}

func (handler *Handler[Conn]) Rotate(builder builder[Conn]) {
	newConnection, err := builder.BuildConnection()
	if err != nil {
		log.Println("err")
	}

	oldConnection := handler.slot.Swap(&slot[Conn]{connection: newConnection})

	go closeGracefully(oldConnection.connection)
}

func closeGracefully[Conn any](conn Connection[Conn]) {
	conn.Wait()
	err := conn.Close()
	if err != nil {
		fmt.Println(err)
	}
}
