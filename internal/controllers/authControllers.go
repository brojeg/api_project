package controllers

import (
	"context"
	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	"encoding/json"
	"net/http"
)

var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	acc := &account.Account{}
	err := json.NewDecoder(r.Body).Decode(acc)
	if err != nil || acc.Login == "" || acc.Password == "" {
		// server.Respond(w, server.Message("Invalid request", 400))
		server.RespondWithMessage(w, 400, "Invalid request")
	}
	resp := account.Login(acc.Login, acc.Password)
	w.Header().Add("Authorization", resp.Message.(string))
	// server.Respond(w, resp)
	server.RespondWithMessage(w, resp.ServerCode, resp.Message)
}
var JwtAuthenticationMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := auth.ValidateToken(r)
		if resp.ServerCode != 200 {
			server.RespondWithMessage(w, resp.ServerCode, resp.Message)
			return
		}

		ctx := context.WithValue(r.Context(), auth.ContextUserKey, resp.Message.(uint))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
