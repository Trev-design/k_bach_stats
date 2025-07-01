package mailer

import (
	"log"

	email "github.com/xhit/go-simple-mail/v2"
)

func (srv *Mailhost) ListenForEmails() {
	for {
		select {
		case msg := <-srv.MessageChannel:
			go srv.sendUserValidationEmail(&msg)
		case err := <-srv.ErrorChannel:
			log.Println(err)

		case <-srv.DoneChannel:
			return
		}
	}
}

func (srv *Mailhost) sendUserValidationEmail(request *MessageRequest) {
	srv.Wait.Add(1)
	defer srv.Wait.Done()

	formatedMessage, err := request.Message.createEmailTemplate()
	if err != nil {
		log.Println(err.Error())
		srv.ErrorChannel <- err
		return
	}

	request.Message.From = srv.FromAddress
	request.Message.FromName = srv.FromName

	client, err := srv.setupClient()
	if err != nil {
		log.Println(err.Error())
		srv.ErrorChannel <- err
		return
	}

	mail := email.NewMSG()

	mail.SetFrom(request.Message.From).
		AddTo(request.Message.To).
		SetSubject(request.Message.Subject).
		AddAlternative(email.TextHTML, formatedMessage)

	if err := mail.Send(client); err != nil {
		srv.ErrorChannel <- err
		return
	}
}
