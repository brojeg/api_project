package controllers

import (
	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	"encoding/json"
	"net/http"
)

var CreateAccountWithBalance = func(w http.ResponseWriter, r *http.Request) {
	account := &account.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		server.Respond(w, server.Message("Invalid request", 400))
	}
	resp := account.Create()
	if resp.ServerCode == 200 {
		balance.Create(account.ID)
		w.Header().Add("Authorization", account.Token)
	}
	server.Respond(w, resp)
}
