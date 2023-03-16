package controllers

import (
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	log "diploma/go-musthave-diploma-tpl/pkg/logger"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"go.uber.org/zap"
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

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(server.LimitMiddleware)
	router.Use(server.LoggingMiddleware)
	router.Post("/api/user/register", CreateAccountWithBalance)
	router.Post("/api/user/login", Authenticate)
	router.Group(func(r chi.Router) {
		r.Use(JwtAuthenticationMiddleware)
		r.Post("/api/user/orders", CreateOrder)
		r.Get("/api/user/orders", GetOrders)
		r.Get("/api/user/balance", GetBalance)
		r.Post("/api/user/balance/withdraw", WithdrawFromBalance)
		r.Get("/api/user/withdrawals", GetBalancHistory)
	})
	return router
}
