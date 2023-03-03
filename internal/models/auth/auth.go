package models

import (
	"context"

	"net/http"
	"os"

	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"

	"github.com/dgrijalva/jwt-go"
)

type contextKey string

var (
	ContextUserKey = contextKey("user")
)

func GetUserFromContext(ctx context.Context) (uint, bool) {
	caller, ok := ctx.Value(ContextUserKey).(uint)
	return caller, ok
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
		tokenHeader := r.Header.Get("Authorization")
		tk := &account.Token{}
		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})
		if err != nil {
			response := server.Message("Malformed authentication token", 401)
			server.Respond(w, response)
			return
		}
		if !token.Valid {
			response := server.Message("Token is not valid.", 400)
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			server.Respond(w, response)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserKey, tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
