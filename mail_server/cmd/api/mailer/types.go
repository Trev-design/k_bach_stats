package mailer

import (
	email "github.com/xhit/go-simple-mail/v2"
)

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
	Message       *Message
}

type MessagePayload struct {
	Kind             string `json:"kind"`
	Email            string `json:"email"`
	ValidationNumber string `json:"verify"`
	Name             string `json:"name"`
}

type ValidationMessage struct {
	ValidationNumber string
	Kind             string
	Name             string
}

func (host *Mailhost) getEncryption() email.Encryption {
	switch host.Encryption {
	case "tls":
		return email.EncryptionSTARTTLS

	case "ssl":
		return email.EncryptionSSLTLS

	case "none":
		return email.EncryptionNone

	default:
		return email.EncryptionSTARTTLS
	}
}
