package middleware

import (
	"fmt"
	"net/http"
	"net/url"
)

func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		fmt.Println(origin)

		if origin == "" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		parsedURL, err := url.Parse(origin)
		if err != nil {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		host := parsedURL.Host
		fmt.Println(host)

		if host != "localhost:5173" {
			writer.WriteHeader(http.StatusBadRequest)
			return
		}

		writer.Header().Set("Access-Control-Allow-Origin", origin)

		if request.Method == "OPTIONS" {
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization")
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(writer, request)
	})
}
