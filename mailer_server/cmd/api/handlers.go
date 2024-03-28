package api

import (
	"fmt"
	"mailer-server/cmd/email"
	"net/http"
)

func (server *app) handleSendVerify() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		server.makeResponse(writer, request, "")
	}
}

func (server *app) handleSendChangePass() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		server.makeResponse(writer, request, "new-password@company.com")
	}
}

func (server *app) makeResponse(writer http.ResponseWriter, request *http.Request, emailAddr string) {
	messageRequest := new(email.VerifyMessageRequest)
	if err := readJSON(request, messageRequest); err != nil {
		writeJSON(writer, http.StatusBadRequest, errorMessage{Message: err.Error()})
		return
	}

	message := email.Message{
		To:      messageRequest.To,
		Subject: messageRequest.Subject,
		Data:    messageRequest.Data,
		From:    emailAddr,
	}

	fmt.Println("sending email")

	server.mailer.Wait.Add(1)
	server.mailer.MailerChannel <- message

	writeJSON(writer, 200, struct {
		Message string `json:"message"`
	}{Message: "OK"})
}
