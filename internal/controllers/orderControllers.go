package controllers

import (
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	order "diploma/go-musthave-diploma-tpl/internal/models/order"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	math "diploma/go-musthave-diploma-tpl/pkg/math"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var CreateOrder = func(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Could not get user from context")
		server.Respond(w, server.Message("Could not get user from context", 500))
	}
	rawOrderNumber, errNumber := getRawOrderNumber(r.Body)
	if errNumber != nil {
		server.Respond(w, server.Message("Bad order number format", 422))
	}
	if !math.IsLuhnValid(rawOrderNumber) {
		server.Respond(w, server.Message("Bad order number format", 422))
	} else {
		order := &order.Order{Status: "NEW", UserID: user, UploadedAt: time.Now(), Number: rawOrderNumber}
		resp := order.Create()
		server.Respond(w, resp)
	}

}

var GetOrders = func(w http.ResponseWriter, r *http.Request) {

	resp := server.Response{}
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.Respond(w, server.Message("Could not get user from context", 500))
	}
	data := order.GetOrders(user)
	resp = server.Message("success", 200)
	if len(data) == 0 {
		resp = server.Message("No orders to display", 204)
	}
	resp.Message = data
	server.Respond(w, resp)
}

var GetOrder = func(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	number := param["number"]

	order := order.GetOrderByNumber(number)
	resp := server.Message("success", 200)
	if order == nil {
		resp = server.Message("No orders to display", 204)
	}
	resp.Message = order
	server.Respond(w, resp)
}

func getRawOrderNumber(body io.Reader) (string, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
