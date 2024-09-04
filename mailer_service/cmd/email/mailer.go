package email

import (
	"fmt"
	"log"

	mailer "github.com/xhit/go-simple-mail/v2"
)

func (mail *MailHost) SendUserValidationEmail(request *MessageRequest) {
	defer mail.Wait.Done()

	log.Println("start preparing message")

	formatedMessage, err := request.MSG.createValidationEmailTemplate()
	if err != nil {
		log.Println(err.Error())
		mail.ErrorChannel <- err
		return
	}

	log.Println("made message")

	request.MSG.From = mail.FromAddress
	request.MSG.FromName = mail.FromName

	smtpClient, err := mail.setupClient()
	if err != nil {
		log.Println(err.Error())
		mail.ErrorChannel <- err
		return
	}

	email := mailer.NewMSG()

	email.
		SetFrom(request.MSG.From).
		AddTo(request.MSG.To).
		SetSubject(request.MSG.Subject).
		AddAlternative(mailer.TextHTML, formatedMessage)

	if err := email.Send(smtpClient); err != nil {
		mail.ErrorChannel <- fmt.Errorf("could not send data. error: %v", err)
		return
	}
}
