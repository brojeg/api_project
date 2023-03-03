package controllers

import (
	m "diploma/go-musthave-diploma-tpl/internal/models"

	"encoding/json"
	"net/http"
	"time"
)

var WithdrawFromBalance = func(w http.ResponseWriter, r *http.Request) {
	user, ok := m.GetUserFromContext(r.Context())
	if !ok {
		m.Respond(w, m.Message("Could not get user from context", 500))
	}
	withdraw := &m.BalanceHistory{ProcessedAt: time.Now(), UserID: user}
	err := json.NewDecoder(r.Body).Decode(withdraw)
	if err != nil {
		m.Respond(w, m.Message("Error while decoding request body", 500))
		return
	}
	currentBalance := m.GetBalance(user)
	resp := currentBalance.Withdraw(withdraw.Sum)
	withdraw.Save()
	m.Respond(w, resp)

}

var GetBalancHistoryFor = func(w http.ResponseWriter, r *http.Request) {
	resp := m.Response{}
	user, ok := m.GetUserFromContext(r.Context())
	if !ok {
		m.Respond(w, m.Message("Could not get user from context", 500))
	}
	data := m.GetBalanceHistory(user)
	resp = m.Message("success", 200)
	resp.Message = data
	m.Respond(w, resp)
}

var GetBalanceFor = func(w http.ResponseWriter, r *http.Request) {
	resp := m.Response{}
	user, ok := m.GetUserFromContext(r.Context())
	if !ok {
		m.Respond(w, m.Message("Could not get user from context", 500))
	}
	data := m.GetBalance(user)
	resp = m.Message("success", 200)
	resp.Message = data
	m.Respond(w, resp)
}
