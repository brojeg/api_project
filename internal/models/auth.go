package models

import (
	"context"
	u "diploma/go-musthave-diploma-tpl/internal/utils"
	"net/http"
	"os"
	"strings"

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

		notAuth := []string{"/api/user/register", "/api/user/login", "/api/orders/"}
		requestPath := r.URL.Path

		for _, value := range notAuth {

			if value == requestPath || strings.Contains(requestPath, "/api/orders/") {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")

		// if tokenHeader == "" {
		// 	response := u.Message(false, "Missing auth token", 401)
		// 	u.Respond(w, response)
		// 	return
		// }

		// splitted := strings.Split(tokenHeader, " ")
		// if len(splitted) != 2 {
		// 	response := u.Message(false, "Invalid/Malformed auth token", 401)
		// 	w.WriteHeader(http.StatusForbidden)
		// 	w.Header().Add("Content-Type", "application/json")
		// 	u.Respond(w, response)
		// 	return
		// }

		// tokenPart := splitted[1]
		tk := &Token{}

		// token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte(os.Getenv("token_password")), nil
		// })
		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil {
			response := u.Message(false, "Malformed authentication token", 401)

			u.Respond(w, response)
			return
		}

		if !token.Valid {
			response := u.Message(false, "Token is not valid.", 400)
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		// ctx := context.WithValue(r.Context(), "user", tk.UserId)
		ctx := context.WithValue(r.Context(), ContextUserKey, tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
