package api

import (
	"mailer-server/cmd/email"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (server *app) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Get("", func(writer http.ResponseWriter, request *http.Request) {
		mail := email.Mail{
			Domain:       "localhost",
			Host:         "localhost",
			Port:         1025,
			Encryption:   "none",
			FromAddress:  "info@company.com",
			FromName:     "info",
			ErrorChannel: make(chan error),
		}

		message := email.Message{
			To:      "me@here.com",
			Subject: "Test Email",
			Data:    "Hello World",
		}

		mail.SendEmail(message, make(chan error))
	})

	mux.Post("send_verify_email", server.handleSendVerify())
	return mux
}
