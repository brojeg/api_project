package controllers

import (
	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	"encoding/json"
	"net/http"
)

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	acc := &account.Account{}
	err := json.NewDecoder(r.Body).Decode(acc)
	if err != nil || acc.Login == "" || acc.Password == "" {
		server.Respond(w, server.Message("Invalid request", 400))
		return
	}
	resp := account.Login(acc.Login, acc.Password)
	w.Header().Add("Authorization", resp.Message.(string))
	server.Respond(w, resp)
}
