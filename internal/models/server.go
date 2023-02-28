package models

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

func Respond(w http.ResponseWriter, data Response) {
	w.Header().Add("Content-Type", "application/json")
	switch data.ServerCode {
	case 200:
		w.WriteHeader(http.StatusOK)
	case 202:
		w.WriteHeader(http.StatusAccepted)
	case 204:
		w.WriteHeader(http.StatusNoContent)
	case 400:
		w.WriteHeader(http.StatusBadRequest)
	case 401:
		w.WriteHeader(http.StatusUnauthorized)
	case 402:
		w.WriteHeader(http.StatusPaymentRequired)
	case 409:
		w.WriteHeader(http.StatusConflict)
	case 429:
		w.WriteHeader(http.StatusTooManyRequests)
	case 422:
		w.WriteHeader(http.StatusUnprocessableEntity)
	case 500:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
	json.NewEncoder(w).Encode(data.Message)
}
