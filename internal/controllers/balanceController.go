package controllers

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"encoding/json"
	"net/http"
	"time"
)

var WithdrawFromBalance = func(w http.ResponseWriter, r *http.Request) {
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		u.Respond(w, u.Message(false, "Could not get user from context", 500))
	}
	// user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	withdraw := &models.BalanceHistory{ProcessedAt: time.Now(), UserID: user}

	err := json.NewDecoder(r.Body).Decode(withdraw)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body", 500))
		return
	}

	// orderNumber := withdraw.Order
	// intSlice := u.StringToIntSlice(orderNumber)
	// if !u.IsLuhnValid(intSlice) {

	// 	u.Respond(w, u.Message(false, "Bad order number format", 422))
	// }
	currentBalance := models.GetBalance(user)

	resp := currentBalance.Withdraw(withdraw.Sum)
	withdraw.Save()
	u.Respond(w, resp)

}

var GetBalancHistoryFor = func(w http.ResponseWriter, r *http.Request) {

	var resp u.Response
	// id := r.Context().Value("user").(uint)
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		u.Respond(w, u.Message(false, "Could not get user from context", 500))
	}
	data := models.GetBalanceHistory(user)
	resp = u.Message(true, "success", 200)
	resp.Message = data
	u.Respond(w, resp)
}

var GetBalanceFor = func(w http.ResponseWriter, r *http.Request) {

	var resp u.Response
	// id := r.Context().Value("user").(uint)
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		u.Respond(w, u.Message(false, "Could not get user from context", 500))
	}
	data := models.GetBalance(user)
	resp = u.Message(true, "success", 200)
	resp.Message = data
	u.Respond(w, resp)
}
