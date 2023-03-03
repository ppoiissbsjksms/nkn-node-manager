package models

//import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
//)

import _ "github.com/mattn/go-sqlite3"


var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open("sqlite3", "wallet.db")

	if err != nil {
		panic("Failed to connect to database!")
	}

	database.AutoMigrate(&Wallet{})

	DB = database
}
