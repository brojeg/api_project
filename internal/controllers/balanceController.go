package controllers

import (
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	balanceHistory "diploma/go-musthave-diploma-tpl/internal/models/balanceHistory"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	math "diploma/go-musthave-diploma-tpl/pkg/math"

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
	}

	if !math.IsLuhnValid(withdraw.Order) {
		server.Respond(w, server.Message("Bad order number format", 422))
	}
	currentBalance := balance.Get(user)
	resp := currentBalance.Withdraw(withdraw.Sum)
	if resp.ServerCode == 200 {
		withdraw.Save()
	}
	server.Respond(w, resp)

}

var GetBalancHistory = func(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.Respond(w, server.Message("Could not get user from context", 500))
	}
	data := balanceHistory.GetBalanceHistory(user)
	resp := server.Response{Message: data, ServerCode: 200}
	server.Respond(w, resp)
}

var GetBalance = func(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.Respond(w, server.Message("Could not get user from context", 500))
	}
	data := balance.Get(user)
	resp := server.Response{Message: data, ServerCode: 200}
	server.Respond(w, resp)
}
