package controllers

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	mod "diploma/go-musthave-diploma-tpl/internal/models"
	math "diploma/go-musthave-diploma-tpl/pkg/math"
	"io"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var CreateOrder = func(w http.ResponseWriter, r *http.Request) {
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Could not get user from context")
		mod.Respond(w, mod.Message("Could not get user from context", 500))
	}
	rawOrderNumber, errNumber := getRawOrderNumber(r.Body)
	if errNumber != nil {
		mod.Respond(w, mod.Message("Bad order number format", 422))
	}
	if !math.IsLuhnValid(rawOrderNumber) {
		mod.Respond(w, mod.Message("Bad order number format", 422))
	} else {
		order := &models.Order{Status: "NEW", UserID: user, UploadedAt: time.Now(), Number: rawOrderNumber}
		resp := order.Create()
		mod.Respond(w, resp)
	}

}

var GetOrdersFor = func(w http.ResponseWriter, r *http.Request) {

	resp := mod.Response{}
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		mod.Respond(w, mod.Message("Could not get user from context", 500))
	}
	data := models.GetOrders(user)
	resp = mod.Message("success", 200)
	if len(data) == 0 {
		resp = mod.Message("No orders to display", 204)
	}
	resp.Message = data
	mod.Respond(w, resp)
}

var GetOrder = func(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	number := param["number"]

	order := models.GetOrderByNumber(number)
	resp := mod.Message("success", 200)
	if order == nil {
		resp = mod.Message("No orders to display", 204)
	}
	resp.Message = order
	mod.Respond(w, resp)
}

func getRawOrderNumber(body io.Reader) (string, error) {
	b, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
