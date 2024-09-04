package api

import "log"

func (application *app) ListenForEmails() {
	for {
		select {
		case message := <-application.mailerChannel:
			go application.mailHost.SendUserValidationEmail(&message)

		case err := <-application.mailHost.ErrorChannel:
			log.Println(err.Error())

		case <-application.mailHost.DoneChannel:
			return
		}
	}
}
