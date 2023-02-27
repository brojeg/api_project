package controllers

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"encoding/json"
	"net/http"
	"time"
)

var WithdrawFromBalance = func(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(uint) //Grab the id of the user that send the request
	withdraw := &models.BalanceHistory{Processed_at: time.Now(), UserId: user}

	err := json.NewDecoder(r.Body).Decode(withdraw)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body", 500))
		return
	}

	orderNumber := withdraw.Order
	intSlice := u.StringToIntSlice(orderNumber)
	if !u.IsLuhnValid(intSlice) {

		u.Respond(w, u.Message(false, "Bad order number format", 422))
	}
	currentBalance := models.GetBalance(user)

	resp := currentBalance.Withdraw(withdraw.Sum)
	withdraw.Save()
	u.Respond(w, resp)

}

var GetBalancHistoryFor = func(w http.ResponseWriter, r *http.Request) {

	var resp u.Response
	id := r.Context().Value("user").(uint)
	data := models.GetBalanceHistory(id)
	resp = u.Message(true, "success", 200)
	resp.Message = data
	u.Respond(w, resp)
}

var GetBalanceFor = func(w http.ResponseWriter, r *http.Request) {

	var resp u.Response
	id := r.Context().Value("user").(uint)
	data := models.GetBalance(id)
	resp = u.Message(true, "success", 200)
	resp.Message = data
	u.Respond(w, resp)
}
