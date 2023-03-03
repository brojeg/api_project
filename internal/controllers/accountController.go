package controllers

import (
	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &account.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		server.Respond(w, server.Message("Invalid request", 400))
		return
	}
	resp := account.Create()
	// account = resp.Message(*account)
	balance := &balance.Balance{UserID: account.ID}
	balance.Save()
	w.Header().Add("Authorization", account.Token)
	server.Respond(w, resp)
}
