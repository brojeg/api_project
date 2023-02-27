package main

import (
	"context"
	"diploma/go-musthave-diploma-tpl/config"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
	"diploma/go-musthave-diploma-tpl/internal/models"
)

func main() {

	// router := mux.NewRouter()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config := config.Init()
	// router.HandleFunc("/api/user/register", controllers.CreateAccount).Methods("POST")
	// router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	// router.HandleFunc("/api/user/orders", controllers.CreateOrder).Methods("POST")
	// router.HandleFunc("/api/user/orders", controllers.GetOrdersFor).Methods("GET")
	// router.HandleFunc("/api/orders/{number}", controllers.GetOrder).Methods("GET")

	// router.HandleFunc("/api/user/balance", controllers.GetBalanceFor).Methods("GET")
	// router.HandleFunc("/api/user/balance/withdraw", controllers.WithdrawFromBalance).Methods("POST")
	// router.HandleFunc("/api/user/withdrawals", controllers.GetBalancHistoryFor).Methods("GET")

	// router.Use(model.LimitMiddleware)
	// router.Use(model.JwtAuthentication)

	// port := os.Getenv("PORT")
	// if port == "" {
	// 	port = "8000" //localhost
	// }

	// fmt.Println(port)

	// go func() {
	// 	err := http.ListenAndServe(":"+port, router)
	// 	if err != nil {
	// 		fmt.Print(err)
	// 	}
	// }()
	models.DatabaseInit(config.Database)
	go models.ApplyAccruals(ctx, config.Interval)
	controllers.NewHTTPServer(config.ServerPort)

}
