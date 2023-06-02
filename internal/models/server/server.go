package model

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message    interface{}
	ServerCode int
}

func Message(message string, serverCode int) Response {
	return Response{Message: message, ServerCode: serverCode}
}

func RespondWithMessage(w http.ResponseWriter, statusCode int, message interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}
