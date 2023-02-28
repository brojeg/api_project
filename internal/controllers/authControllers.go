package controllers

import (
	mod "diploma/go-musthave-diploma-tpl/internal/models"
	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &mod.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		mod.Respond(w, mod.Message("Invalid request", 400))
		return
	}
	resp := account.Create()
	account = resp.Message.(*mod.Account)
	balance := &mod.Balance{UserID: account.ID}
	balance.Save()
	w.Header().Add("Authorization", account.Token)
	mod.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &mod.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil || account.Login == "" || account.Password == "" {
		mod.Respond(w, mod.Message("Invalid request", 400))
		return
	}
	resp := mod.Login(account.Login, account.Password)
	w.Header().Add("Authorization", resp.Message.(string))
	mod.Respond(w, resp)
}
