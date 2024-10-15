package middleware

import (
	"log"
	"net/http"
	"strings"
	"user_manager/database"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Name    string `json:"name"`
	Account string `json:"entity"`
	AboType string `json:"abo"`
	Session string `json:"session_id"`
	jwt.RegisteredClaims
}

func Auth(sessionHandler database.SessionHandler, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Authorization")

		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("no auth header")
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			log.Println("invalid auth header")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if headerParts[0] != "Bearer" {
			log.Println("invalid auth header")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := sessionHandler.CheckSession(headerParts[1]); err != nil {
			log.Println("invalid auth header")
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
