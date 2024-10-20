package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"user_manager/database"
)

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
		fmt.Println(headerParts)
		fmt.Printf("the length of the header array is: %d\n", len(headerParts))
		if len(headerParts) != 2 {
			log.Println("invalid auth header")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Printf("the head of the header is : %s\n", headerParts[0])
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
