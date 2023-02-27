package controllers

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"encoding/json"
	"net/http"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		u.Respond(w, u.Message(false, "Invalid request", 400))
		return
	}

	resp, token := account.Create() //Create account
	w.Header().Add("Authorization", token)
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil || account.Login == "" || account.Password == "" {
		u.Respond(w, u.Message(false, "Invalid request", 400))
		return
	}

	resp := models.Login(account.Login, account.Password)
	u.Respond(w, resp)
}
