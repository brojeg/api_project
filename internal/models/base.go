package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func DatabaseInit(connectionString string) {

	// connectionString = "host=localhost user=test_user dbname=test sslmode=disable password=111"
	conn, err := gorm.Open("postgres", connectionString)
	if err != nil {
		logger.Errorf("Error is %e \n Connection string is %s", err, connectionString)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Order{}, &Balance{}, &BalanceHistory{})
}

func GetDB() *gorm.DB {
	return db
}
