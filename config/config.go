package config

import (
	account "diploma/go-musthave-diploma-tpl/internal/models/account"
	auth "diploma/go-musthave-diploma-tpl/internal/models/auth"
	balance "diploma/go-musthave-diploma-tpl/internal/models/balance"
	balanceHistory "diploma/go-musthave-diploma-tpl/internal/models/balanceHistory"
	db "diploma/go-musthave-diploma-tpl/internal/models/database"
	order "diploma/go-musthave-diploma-tpl/internal/models/order"
	server "diploma/go-musthave-diploma-tpl/internal/models/server"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type ServerConfig struct {
	ServerPort     string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	Interval       string `env:"INTERVAL" envDefault:"5s"`
	Database       string `env:"DATABASE_URI"`
	Accrual        string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"127.0.0.1:8081"`
	JWTPassword    string `env:"JWT_PASSWORD"`
	ExpirationTime int    `env:"EXPIRATION_TIME" envDefault:"15"`
	ServerLog      string `env:"SERVER_LOG"`
}

func Init() ServerConfig {

	godotenv.Load(".env")

	var envCfg ServerConfig
	err := env.Parse(&envCfg)

	_, envAdddressExists := os.LookupEnv("RUN_ADDRESS")
	_, envDBExists := os.LookupEnv("DATABASE_URI")
	_, envAccrualExists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	_, envIntervalExists := os.LookupEnv("INTERVAL")
	_, envJWTPAsswordExists := os.LookupEnv("JWT_PASSWORD")
	_, envExpirationTimeExists := os.LookupEnv("EXPIRATION_TIME")

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
	flag.Func("d", "Postgres connection string (No default value)", func(flagValue string) error {
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
	flag.Func("p", "Interval for the accrual system check (default 5s)", func(flagValue string) error {
		if envJWTPAsswordExists {
			return nil
		}
		envCfg.JWTPassword = flagValue
		return nil
	})
	flag.Func("t", "TTL for JWT token (default 15m", func(flagValue string) error {
		if envExpirationTimeExists {
			return nil
		}
		intVar, err := strconv.Atoi(flagValue)
		if err != nil {
			return err
		}
		envCfg.ExpirationTime = intVar
		return nil
	})
	flag.Parse()

	db.InitDBConnectionString(envCfg.Database)
	order.InitAccrualURL(envCfg.Accrual)
	auth.InitJWTPassword(envCfg.JWTPassword, envCfg.ExpirationTime)
	server.SetServerLogPath(envCfg.ServerLog)
	order.CreateTable()
	account.CreateTable()
	balance.CreateTable()
	balanceHistory.CreateTable()

	return envCfg
}
