package config

import (
	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	balanceHistory "diploma/go-musthave-diploma-tpl/internal/models/balanceHistory"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	order "diploma/go-musthave-diploma-tpl/internal/models/order"
	"flag"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerPort string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	Interval   string `env:"INTERVAL" envDefault:"5s"`
	Database   string `env:"DATABASE_URI"`
	Accrual    string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Init() ServerConfig {
	var envCfg ServerConfig
	err := env.Parse(&envCfg)

	_, envAdddressExists := os.LookupEnv("RUN_ADDRESS")
	_, envDBExists := os.LookupEnv("DATABASE_URI")
	_, envAccrualExists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	_, envIntervalExists := os.LookupEnv("Interval")
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
	flag.Func("a", "Server address (default localhost:8080)", func(flagValue string) error {
		if envAdddressExists {
			return nil
		}
		envCfg.ServerPort = flagValue
		return nil
	})
	flag.Func("d", "Posgres connection string (No default value)", func(flagValue string) error {
		if envDBExists {
			return nil
		}
		envCfg.Database = flagValue

		return nil
	})
	flag.Func("r", "ACCRUAL SYSTEM ADDRESS (No default value)", func(flagValue string) error {
		if envAccrualExists {
			return nil
		}
		envCfg.Accrual = flagValue
		return nil
	})
	flag.Func("i", "Interval for the accrual system check (default 5s)", func(flagValue string) error {
		if envIntervalExists {
			return nil
		}
		envCfg.Interval = flagValue
		return nil
	})
	flag.Parse()

	db.InitDBConnectionString(envCfg.Database)
	order.InitAccrualURL(envCfg.Accrual)
	order.CreteTable()
	account.CreteTable()
	balance.CreteTable()
	balanceHistory.CreteTable()

	return envCfg
}
