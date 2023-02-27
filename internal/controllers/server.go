package controllers

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func NewHTTPServer(port string) {

	router := NewRouter()

	if err := http.ListenAndServe(port, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("Cannot start http.ListenAndServe. Error is: /n %e", err)
	} else {
		fmt.Printf("application stopped gracefully")
	}

}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")

	router.HandleFunc("/api/user/orders", CreateOrder).Methods("POST")
	router.HandleFunc("/api/user/orders", GetOrdersFor).Methods("GET")
	router.HandleFunc("/api/orders/{number}", GetOrder).Methods("GET")

	router.HandleFunc("/api/user/balance", GetBalanceFor).Methods("GET")
	router.HandleFunc("/api/user/balance/withdraw", WithdrawFromBalance).Methods("POST")
	router.HandleFunc("/api/user/withdrawals", GetBalancHistoryFor).Methods("GET")
	router.Use(models.LimitMiddleware)
	router.Use(models.JwtAuthentication)

	return router

}
