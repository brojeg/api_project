package models

import (
	"fmt"

	log "diploma/go-musthave-diploma-tpl/pkg/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"go.uber.org/zap"
)

var connectionString string
var logger *zap.SugaredLogger = log.Init()

func InitDBConnectionString(conn string) {
	connectionString = conn
}

func Get() *gorm.DB {
	conn, err := gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Println(connectionString)
		logger.Errorf("Error is %e \n Connection string is %s", err, connectionString)
	}
	return conn
}
