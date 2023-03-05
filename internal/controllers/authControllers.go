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
		server.Respond(w, server.Message("Invalid request", 400))
	}
	resp := account.Login(acc.Login, acc.Password)
	w.Header().Add("Authorization", resp.Message.(string))
	server.Respond(w, resp)
}
var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"/api/user/register", "/api/user/login"}
		requestPath := r.URL.Path
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		resp := auth.ValidateToken(r)
		if resp.ServerCode != 200 {
			server.Respond(w, resp)
		}
		ctx := context.WithValue(r.Context(), auth.ContextUserKey, resp.Message.(uint))
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
