package models

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func DatabaseInit(connectionString string) {
	time.Sleep(2 * time.Second)
	// connectionString = "host=localhost user=test_user dbname=test sslmode=disable password=111"
	conn, err := gorm.Open("postgres", connectionString)
	if err != nil {
		logger.Error(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Order{}, &Balance{}, &BalanceHistory{})
}

func GetDB() *gorm.DB {
	return db
}
