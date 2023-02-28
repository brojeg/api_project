package controllers

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message("Invalid request", 400))
		return
	}
	resp := account.Create()
	account = resp.Message.(*models.Account)
	balance := &models.Balance{UserID: account.ID}
	balance.Save()
	w.Header().Add("Authorization", account.Token)
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil || account.Login == "" || account.Password == "" {
		u.Respond(w, u.Message("Invalid request", 400))
		return
	}
	resp := models.Login(account.Login, account.Password)
	w.Header().Add("Authorization", resp.Message.(string))
	u.Respond(w, resp)
}
