package api

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message string `json:"message"`
}

func writeJSON(writer http.ResponseWriter, status int, body any) error {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status)
	return json.NewEncoder(writer).Encode(body)
}

func readJSON(request *http.Request, body any) error {
	return json.NewDecoder(request.Body).Decode(body)
}
