package email

import (
	"sync"

	mailer "github.com/xhit/go-simple-mail/v2"
)

type MailHost struct {
	Domain       string
	Host         string
	Port         int
	UserName     string
	Password     string
	Encryption   string
	FromAddress  string
	FromName     string
	Wait         *sync.WaitGroup
	ErrorChannel chan error
	DoneChannel  chan bool
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Payload     any
}

type MessageRequest struct {
	CorrelationID string
	UserID        string
	Email         string
	MSG           *Message
}

type ValidationMessagePayload struct {
	Kind             string `json:"kind"`
	Email            string `json:"email"`
	ValidationNumber string `json:"verify"`
	Name             string `json:"name"`
	UserId           string `json:"user_id"`
}

type ValidationMessage struct {
	ValidationNumber string
	Name             string
	Kind             string
}

func (host *MailHost) getEncryption() mailer.Encryption {
	switch host.Encryption {
	case "tls":
		return mailer.EncryptionSTARTTLS

	case "ssl":
		return mailer.EncryptionSSLTLS

	case "none":
		return mailer.EncryptionNone

	default:
		return mailer.EncryptionSTARTTLS
	}
}

func NewValidationMailer() *MailHost {
	return &MailHost{
		Domain:       "localhost",
		Host:         "localhost",
		Port:         1025,
		Encryption:   "none",
		FromAddress:  "support@kbach.com",
		FromName:     "support",
		ErrorChannel: make(chan error),
		DoneChannel:  make(chan bool),
		Wait:         &sync.WaitGroup{},
	}
}
