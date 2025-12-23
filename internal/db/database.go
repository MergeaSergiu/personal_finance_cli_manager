package db

import (
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		if err := os.MkdirAll("data", os.ModePerm); err != nil {
			log.Fatalf("Failed to create data folder: %v", err)
		}
	}

	db, err := gorm.Open(sqlite.Open("data/app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	DB = db
}
