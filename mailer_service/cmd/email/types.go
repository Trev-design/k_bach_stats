package email

import (
	"sync"

	mailer "github.com/xhit/go-simple-mail/v2"
)

type MailHost struct {
	Domain        string
	Host          string
	Port          int
	UserName      string
	Password      string
	Encryption    string
	FromAddress   string
	FromName      string
	Wait          *sync.WaitGroup
	MailerChannel chan Message
	ErrorChannel  chan error
	DoneChannel   chan bool
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Payload     any
}

type ValidationMessage struct {
	Kind             string
	Email            string
	ValidationNumber string
	Name             string
}

func (mail *MailHost) getEncryption() mailer.Encryption {
	switch mail.Encryption {
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
