package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func DatabaseInit(connectionString string) {

	conn, err := gorm.Open("postgres", connectionString)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Order{}, &Balance{}, &BalanceHistory{})
}

func GetDB() *gorm.DB {
	return db
}
