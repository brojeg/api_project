package config

import (
	"diploma/go-musthave-diploma-tpl/internal/models"
	"flag"
	"fmt"
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
			fmt.Println("env value for server exists")
			return nil
		}
		envCfg.Database = flagValue
		fmt.Println("Unsing Flag value for server")
		return nil
	})
	flag.Func("r", "ACCRUAL SYSTEM ADDRESS (No default value)", func(flagValue string) error {
		if envAccrualExists {

			return nil
		}
		envCfg.Accrual = flagValue
		return nil
	})
	flag.Parse()

	models.InitDBConnectionString(envCfg.Database)
	models.InitAccrualURL(envCfg.Accrual)

	return envCfg
}
