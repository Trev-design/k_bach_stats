package email

import (
	"fmt"
	"log"
	"sync"

	mailer "github.com/xhit/go-simple-mail/v2"
)

func NewValidationMailer() *MailHost {
	return &MailHost{
		Domain:        "localhost",
		Host:          "localhost",
		Port:          1025,
		Encryption:    "none",
		FromAddress:   "support@company.com",
		FromName:      "support",
		ErrorChannel:  make(chan error),
		MailerChannel: make(chan Message, 100),
		DoneChannel:   make(chan bool),
		Wait:          &sync.WaitGroup{},
	}
}

func (mail *MailHost) SendUserValidationEmail(message *Message) {
	defer mail.Wait.Done()

	formatedMessage, err := message.createValidationEmailTemplate()
	if err != nil {
		mail.ErrorChannel <- err
		return
	}

	message.From = mail.FromAddress
	message.FromName = mail.FromName

	smtpClient, err := mail.setupClient()
	if err != nil {
		mail.ErrorChannel <- err
	}

	email := mailer.NewMSG()

	email.SetFrom(message.From).AddTo(message.To).SetSubject(message.Subject)
	email.AddAlternative(mailer.TextHTML, formatedMessage)

	if err := email.Send(smtpClient); err != nil {
		mail.ErrorChannel <- fmt.Errorf("could not send email: %v", err)
	}
}

func (mail *MailHost) ListenForEmails() {

	for {
		select {
		case message := <-mail.MailerChannel:
			log.Println("got message: ", message)
			go mail.SendUserValidationEmail(&message)

		case err := <-mail.ErrorChannel:
			log.Println(err.Error())

		case <-mail.DoneChannel:
			return
		}
	}
}
