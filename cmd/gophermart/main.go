package main

import (
	"context"
	"diploma/go-musthave-diploma-tpl/config"
	"diploma/go-musthave-diploma-tpl/internal/controllers"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config.Init()
	config.InitLocalVars()
	config.CreateDBTables()
	go controllers.ApplyAccruals(ctx, config.Param.Interval, config.Param.Accrual)
	controllers.NewHTTPServer(config.Param.ServerPort)

}
