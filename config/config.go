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
	HTTPServer
	ExternalDependency
	ServerAuth
	ServerLog
}
type HTTPServer struct {
	ServerPort string `env:"RUN_ADDRESS" envDefault:"127.0.0.1:8080"`
	Interval   string `env:"INTERVAL" envDefault:"5s"`
}
type ExternalDependency struct {
	Database string `env:"DATABASE_URI"`
	Accrual  string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"127.0.0.1:8081"`
}
type ServerAuth struct {
	JWTPassword    string `env:"JWT_PASSWORD"`
	ExpirationTime int    `env:"EXPIRATION_TIME" envDefault:"15"`
}
type ServerLog struct {
	Log string `env:"SERVER_LOG"`
}

var Param ServerConfig

func InitStartupParameters() {

	godotenv.Load(".env")
	err := env.Parse(&Param)

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
		Param.ServerPort = flagValue
		return nil
	})
	flag.Func("d", "Postgres connection string (No default value)", func(flagValue string) error {
		if envDBExists {
			return nil
		}
		Param.Database = flagValue

		return nil
	})
	flag.Func("r", "ACCRUAL SYSTEM ADDRESS (No default value)", func(flagValue string) error {
		if envAccrualExists {
			return nil
		}
		Param.Accrual = flagValue
		return nil
	})
	flag.Func("i", "Interval for the accrual system check (default 5s)", func(flagValue string) error {
		if envIntervalExists {
			return nil
		}
		Param.Interval = flagValue
		return nil
	})
	flag.Func("p", "Interval for the accrual system check (default 5s)", func(flagValue string) error {
		if envJWTPAsswordExists {
			return nil
		}
		Param.JWTPassword = flagValue
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
		Param.ExpirationTime = intVar
		return nil
	})
	flag.Parse()
}

func InitLocalVars() {
	db.InitDBConnectionString(Param.Database)
	order.InitAccrualURL(Param.Accrual)
	auth.InitJWTPassword(Param.JWTPassword, Param.ExpirationTime)
	server.SetServerLogPath(Param.Log)
}

func CreateDBTables() {
	order.CreateTable()
	account.CreateTable()
	balance.CreateTable()
	balanceHistory.CreateTable()
}
