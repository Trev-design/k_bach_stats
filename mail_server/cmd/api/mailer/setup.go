package mailer

import "sync"

type Mailhost struct {
	Domain         string
	Host           string
	Port           int
	Username       string
	Password       string
	Encryption     string
	FromAddress    string
	FromName       string
	Wait           *sync.WaitGroup
	ErrorChannel   chan error
	DoneChannel    chan struct{}
	MessageChannel <-chan MessageRequest
}

func NewMailHost(mailerChannel <-chan MessageRequest) *Mailhost {
	return &Mailhost{
		Domain:         "localhost",
		Host:           "localhost",
		Port:           1025,
		Encryption:     "none",
		FromAddress:    "support@kbach.com",
		FromName:       "kbach support",
		ErrorChannel:   make(chan error),
		DoneChannel:    make(chan struct{}),
		MessageChannel: mailerChannel,
		Wait:           &sync.WaitGroup{},
	}
}

func (srv *Mailhost) CloseMailhost() {
	srv.DoneChannel <- struct{}{}
	srv.Wait.Wait()
	close(srv.ErrorChannel)
	close(srv.DoneChannel)
}
