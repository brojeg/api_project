package controllers

import (
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	log "diploma/go-musthave-diploma-tpl/pkg/logger"
	"errors"
	"net/http"

	"go.uber.org/zap"

	"github.com/gorilla/mux"
)

var logger *zap.SugaredLogger = log.Init()

func NewHTTPServer(port string) {

	router := NewRouter()

	if err := http.ListenAndServe(port, router); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("Cannot start http.ListenAndServe. Error is: /n %e", err)
	} else {
		logger.Warn("application stopped gracefully")
	}

}

func NewRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", CreateAccountWithBalance).Methods("POST")
	router.HandleFunc("/api/user/login", Authenticate).Methods("POST")
	router.HandleFunc("/api/user/orders", CreateOrder).Methods("POST")
	router.HandleFunc("/api/user/orders", GetOrders).Methods("GET")
	router.HandleFunc("/api/user/balance", GetBalance).Methods("GET")
	router.HandleFunc("/api/user/balance/withdraw", WithdrawFromBalance).Methods("POST")
	router.HandleFunc("/api/user/withdrawals", GetBalancHistory).Methods("GET")
	router.Use(server.LimitMiddleware)
	router.Use(auth.JwtAuthentication)
	router.Use(server.HTTPLogger)

	return router

}
