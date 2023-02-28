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
		u.Respond(w, u.Message("Could not get user from context", 500))
	}
	withdraw := &models.BalanceHistory{ProcessedAt: time.Now(), UserID: user}
	err := json.NewDecoder(r.Body).Decode(withdraw)
	if err != nil {
		u.Respond(w, u.Message("Error while decoding request body", 500))
		return
	}
	currentBalance := models.GetBalance(user)
	resp := currentBalance.Withdraw(withdraw.Sum)
	withdraw.Save()
	u.Respond(w, resp)

}

var GetBalancHistoryFor = func(w http.ResponseWriter, r *http.Request) {
	resp := u.Response{}
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		u.Respond(w, u.Message("Could not get user from context", 500))
	}
	data := models.GetBalanceHistory(user)
	resp = u.Message("success", 200)
	resp.Message = data
	u.Respond(w, resp)
}

var GetBalanceFor = func(w http.ResponseWriter, r *http.Request) {
	resp := u.Response{}
	user, ok := models.GetUserFromContext(r.Context())
	if !ok {
		u.Respond(w, u.Message("Could not get user from context", 500))
	}
	data := models.GetBalance(user)
	resp = u.Message("success", 200)
	resp.Message = data
	u.Respond(w, resp)
}
