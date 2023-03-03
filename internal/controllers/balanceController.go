package controllers

import (
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	balanceHistory "diploma/go-musthave-diploma-tpl/internal/models/balanceHistory"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"

	"encoding/json"
	"net/http"
	"time"
)

var WithdrawFromBalance = func(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.Respond(w, server.Message("Could not get user from context", 500))
	}
	withdraw := &balanceHistory.BalanceHistory{ProcessedAt: time.Now(), UserID: user}
	err := json.NewDecoder(r.Body).Decode(withdraw)
	if err != nil {
		server.Respond(w, server.Message("Error while decoding request body", 500))
		return
	}
	currentBalance := balance.Get(user)
	resp := currentBalance.Withdraw(withdraw.Sum)
	withdraw.Save()
	server.Respond(w, resp)

}

var GetBalancHistoryFor = func(w http.ResponseWriter, r *http.Request) {
	resp := server.Response{}
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.Respond(w, server.Message("Could not get user from context", 500))
	}
	data := balanceHistory.GetBalanceHistory(user)
	resp = server.Message("success", 200)
	resp.Message = data
	server.Respond(w, resp)
}

var GetBalanceFor = func(w http.ResponseWriter, r *http.Request) {
	resp := server.Response{}
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.Respond(w, server.Message("Could not get user from context", 500))
	}
	data := balance.Get(user)
	resp = server.Message("success", 200)
	resp.Message = data
	server.Respond(w, resp)
}
