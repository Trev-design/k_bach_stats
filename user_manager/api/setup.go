package api

import (
	"log"
	"user_manager/database/db"
	redissession "user_manager/database/redis_session"
	"user_manager/rabbitmq"
)

type App struct {
	Session      *redissession.SessionClient
	DBIsntance   *db.Database
	Consumer     rabbitmq.RabbitConsumer
	ErrorChannel chan error
	DoneChannel  chan bool
}

type ServerOptions struct {
	Database string
	Session  string
}

func (options *ServerOptions) Setup() (*App, error) {
	log.Println("init session")
	session, err := redissession.Setup()
	if err != nil {
		return nil, err
	}

	log.Println("init user database")
	db, err := db.Setup(options.Database)
	if err != nil {
		return nil, err
	}

	log.Println("init consumer structure")
	consumer, err := rabbitmq.Setup()
	if err != nil {
		return nil, err
	}

	return &App{
		DBIsntance:   db,
		Session:      session,
		Consumer:     consumer,
		ErrorChannel: make(chan error),
		DoneChannel:  make(chan bool),
	}, nil
}
