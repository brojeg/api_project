package main

import (
	"context"
	"diploma/go-musthave-diploma-tpl/config"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config := config.Init()
	go controllers.ApplyAccruals(ctx, config.Interval, config.Accrual)
	controllers.NewHTTPServer(config.ServerPort)

}
