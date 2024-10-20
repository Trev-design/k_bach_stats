package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
	"user_manager/api"
	"user_manager/database"
)

func InitialAuth(sessionHandler database.SessionHandler, next http.Handler) http.Handler {

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log.Println("welcome in initial auth")
		writer.Header().Add("Vary", "Authorization")

		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("no authheader")
			writer.WriteHeader(http.StatusUnauthorized)
			return
		}

		log.Println("have an auth header")

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			log.Println("invalid authheader")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		log.Println("have good chances for a valid auth header")

		if headerParts[0] != "Bearer" {
			log.Println()
			writer.WriteHeader(http.StatusConflict)
			return
		}

		log.Println("better chances for a valid auth header")

		entity, err := sessionHandler.InitialAuth(headerParts[1])
		if err != nil {
			writer.WriteHeader(http.StatusForbidden)
		}

		log.Println("auth header is valid")

		ctx := context.WithValue(request.Context(), api.ContextKey("entity"), entity)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}
