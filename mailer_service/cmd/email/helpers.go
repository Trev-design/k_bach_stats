package email

import (
	"bytes"
	"context"
	"fmt"
	"mailer_service/cmd/templates"
	"time"

	mailer "github.com/xhit/go-simple-mail/v2"
)

func (message *Message) createValidationEmailTemplate() (string, error) {
	payload, ok := message.Payload.(*ValidationMessage)
	if !ok {
		return "", fmt.Errorf("invalid credentials")
	}

	formatMessage, err := generateValidationTemplate(payload)
	if err != nil {
		return "", err
	}

	return formatMessage, nil
}

func generateValidationTemplate(payload *ValidationMessage) (string, error) {
	var buffer bytes.Buffer

	if err := templates.VerifyUser(
		payload.ValidationNumber,
		payload.Name,
		payload.Kind,
	).Render(
		context.Background(),
		&buffer,
	); err != nil {
		return "", fmt.Errorf("could not generate email. error: %v", err)
	}

	return buffer.String(), nil
}

func (mail *MailHost) setupClient() (*mailer.SMTPClient, error) {
	server := mailer.NewSMTPClient()
	server.Host = mail.Host
	server.Port = mail.Port
	server.Username = mail.UserName
	server.Password = mail.Password
	server.Encryption = mail.getEncryption()
	server.ConnectTimeout = 15 * time.Second
	server.SendTimeout = 15 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		return nil, fmt.Errorf("could not connect to mailer client, error: %v", err)
	}

	return smtpClient, err
}
