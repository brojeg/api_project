package config

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v6"
)

type ServerConfig struct {
	ServerPort string `env:"ADDRESS" envDefault:"127.0.0.1:8080"`
	Interval   string `env:"INTERVAL" envDefault:"5s"`
	Database   string `env:"DATABASE_DSN" envDefault:"postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"`
	Accrual    string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func Init() ServerConfig {

	return getParams()
}

func getParams() ServerConfig {

	var envCfg ServerConfig
	err := env.Parse(&envCfg)

	if err != nil {

		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}

	flag.Func("a", "Server address (default localhost:8080)", func(flagValue string) error {

		_, exists := os.LookupEnv("ADDRESS")

		if exists {
			return nil
		}

		envCfg.ServerPort = flagValue
		return nil
	})
	flag.Func("d", "Posgres connection string (No default value)", func(flagValue string) error {

		_, exists := os.LookupEnv("DATABASE_DSN")

		if exists {
			return nil
		}

		envCfg.Database = flagValue
		return nil
	})
	flag.Func("r", "ACCRUAL SYSTEM ADDRESS (No default value)", func(flagValue string) error {

		_, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")

		if exists {
			return nil
		}

		envCfg.Accrual = flagValue
		return nil
	})

	flag.Parse()

	fmt.Println(envCfg)
	return envCfg
}
