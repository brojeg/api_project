package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var connectionString string

func InitDBConnectionString(conn string) {
	connectionString = conn
}

func getConstring() string {
	return connectionString
}

func GetDB() *gorm.DB {

	conn, err := gorm.Open("postgres", getConstring())
	if err != nil {
		fmt.Println(connectionString)
		logger.Errorf("Error is %e \n Connection string is %s", err, connectionString)
	}

	conn.Debug().AutoMigrate(&Account{}, &Order{}, &Balance{}, &BalanceHistory{})

	return conn
}
