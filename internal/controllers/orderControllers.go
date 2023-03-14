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
		server.RespondWithMessage(w, 500, "Could not get user from context")
	}
	rawOrderNumber, errNumber := getRawOrderNumber(r.Body)
	if errNumber != nil {
		server.RespondWithMessage(w, 422, "Bad order number format")
	}
	if !math.IsLuhnValid(rawOrderNumber) {
		server.RespondWithMessage(w, 422, "Bad order number format")
	} else {
		order := &order.Order{Status: "NEW", UserID: user, UploadedAt: time.Now(), Number: rawOrderNumber}
		resp := order.Create()
		server.RespondWithMessage(w, resp.ServerCode, resp.Message)
	}

}

var GetOrders = func(w http.ResponseWriter, r *http.Request) {

	resp := server.Response{}
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.RespondWithMessage(w, 500, "Could not get user from context")
	}
	data := order.GetOrders(user)
	resp = server.Message("success", 200)
	if len(data) == 0 {
		resp = server.Message("No orders to display", 204)
	}
	resp.Message = data
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)
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
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)
}

func getRawOrderNumber(body io.Reader) (string, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
