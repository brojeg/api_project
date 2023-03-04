package models

import (
	"context"

	"net/http"

	server "diploma/go-musthave-diploma-tpl/internal/models/server"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

var (
	ContextUserKey = contextKey("user")
)
var jwtPassword string

func InitJWTPassword(pass string) {
	jwtPassword = pass
}

type Token struct {
	UserID uint
	jwt.StandardClaims
}

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
		tk := &Token{}
		token, err := jwt.ParseWithClaims(tokenHeader, tk, func(token *jwt.Token) (interface{}, error) {
			// return []byte(os.Getenv("token_password")), nil
			return []byte(jwtPassword), nil
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

func GetToken(id uint) string {
	tk := &Token{UserID: id}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, err := token.SignedString([]byte(jwtPassword))
	if err != nil {
		panic(err)
	}

	return tokenString
}

func EncryptPassword(pass string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func IsPasswordsEqual(existing, new string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(existing), []byte(new))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false
	}
	return true
}
