package controllers

import (
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	balanceHistory "diploma/go-musthave-diploma-tpl/internal/models/balanceHistory"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	math "diploma/go-musthave-diploma-tpl/pkg/math"

	"encoding/json"
	"net/http"
)

var WithdrawFromBalance = func(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.RespondWithMessage(w, 500, "Could not get user from context")
	}

	withdraw := balanceHistory.NewBalance(user)
	err := json.NewDecoder(r.Body).Decode(withdraw)
	if err != nil {
		// server.Respond(w, server.Message("Error while decoding request body", 500))
		server.RespondWithMessage(w, 500, "Error while decoding request body")
	}

	if !math.IsLuhnValid(withdraw.Order) {
		// server.Respond(w, server.Message("Bad order number format", 422))
		server.RespondWithMessage(w, 422, "Bad order number format")
	}
	currentBalance := balance.Get(user)
	resp := currentBalance.Withdraw(withdraw.Sum)
	if resp.ServerCode == 200 {
		withdraw.Save()
	}
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)

}

var GetBalancHistory = func(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.RespondWithMessage(w, 500, "Could not get user from context")
	}
	data := balanceHistory.GetBalanceHistory(user)
	resp := server.Response{Message: data, ServerCode: 200}
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)
}

var GetBalance = func(w http.ResponseWriter, r *http.Request) {
	user, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		server.RespondWithMessage(w, 500, "Could not get user from context")
	}
	data := balance.Get(user)
	resp := server.Response{Message: data, ServerCode: 200}
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)
}
