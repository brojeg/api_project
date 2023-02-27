package main

import (
	"context"
	"diploma/go-musthave-diploma-tpl/config"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
	"diploma/go-musthave-diploma-tpl/internal/models"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config := config.Init()
	models.InitDBConnectionString(config.Database)
	go models.ApplyAccruals(ctx, config.Interval, config.Accrual)
	controllers.NewHTTPServer(config.ServerPort)

}
