package api

import (
	"mailer-server/cmd/email"
	"net/http"
)

func (server *app) handleSendVerify() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		messageRequest := new(email.VerifyMessageRequest)
		if err := readJSON(request, messageRequest); err != nil {
			writeJSON(writer, http.StatusBadRequest, errorMessage{Message: err.Error()})
		}

		message := email.Message{
			To:      messageRequest.To,
			Subject: messageRequest.Subject,
			Data:    messageRequest.Data,
		}

		server.mailer.Wait.Add(1)
		server.mailer.MailerChannel <- message
	}
}
