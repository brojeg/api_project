package controllers

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var CreateOrder = func(w http.ResponseWriter, r *http.Request) {
	// user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		logger.Error("Could not get user from context")
		u.Respond(w, u.Message(false, "Could not get user from context", 500))
	}

	rawOrderNumber := u.GetRawOrderNumber(r.Body)
	if !u.IsLuhnValid(u.StringToIntSlice(rawOrderNumber)) {
		u.Respond(w, u.Message(false, "Bad order number format", 422))
	} else {
		order := &models.Order{Status: "NEW", UserID: user, UploadedAt: time.Now(), Number: rawOrderNumber}
		resp := order.Create()

		u.Respond(w, resp)
	}

}

var GetOrdersFor = func(w http.ResponseWriter, r *http.Request) {

	var resp u.Response
	// id := r.Context().Value("user").(uint)
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		u.Respond(w, u.Message(false, "Could not get user from context", 500))
	}
	data := models.GetOrders(user)
	resp = u.Message(true, "success", 200)
	if len(data) == 0 {
		resp = u.Message(false, "No orders to display", 204)
	}
	resp.Message = data
	u.Respond(w, resp)
}

var GetOrder = func(w http.ResponseWriter, r *http.Request) {
	param := mux.Vars(r)
	number := param["number"]

	order := models.GetOrderByNumber(number)
	resp := u.Message(true, "success", 200)
	if order == nil {
		resp = u.Message(false, "No orders to display", 204)
	}
	resp.Message = order
	u.Respond(w, resp)
}
