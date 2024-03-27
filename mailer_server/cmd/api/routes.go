package api

import (
	"mailer-server/cmd/email"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (server *app) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

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
