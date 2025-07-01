package mailer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mail_server/cmd/api/templates"
	"time"

	email "github.com/xhit/go-simple-mail/v2"
)

func (message *Message) createEmailTemplate() (string, error) {
	payload, ok := message.Payload.(*ValidationMessage)
	if !ok {
		return "", errors.New("invalid credentials")
	}

	return payload.generateTemplate()
}

func (payload *ValidationMessage) generateTemplate() (string, error) {
	var buffer *bytes.Buffer

	if err := templates.VerifyUser(
		payload.ValidationNumber,
		payload.Name,
		payload.Kind,
	).Render(
		context.Background(),
		buffer,
	); err != nil {
		return "", fmt.Errorf("could not generate emails %v", err)
	}

	return buffer.String(), nil
}

func (srv *Mailhost) setupClient() (*email.SMTPClient, error) {
	server := email.NewSMTPClient()
	server.Host = srv.Host
	server.Port = srv.Port
	server.Username = srv.Username
	server.Password = srv.Password
	server.Encryption = srv.getEncryption()
	server.ConnectTimeout = 15 * time.Second
	server.SendTimeout = 15 * time.Second

	return server.Connect()
}
