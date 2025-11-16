package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() {

	if _, err := os.Stat("../sqlite"); os.IsNotExist(err) {
		if err := os.MkdirAll("../sqlite", os.ModePerm); err != nil {
			log.Fatalf("Failed to create db folder: %v", err)
		}
	}

	var err error
	DB, err = sql.Open("sqlite", "../sqlite/app.db") // database file inside db folder
	if err != nil {
		log.Fatal(err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT,
		description TEXT,
		amount REAL,
		category TEXT
	);`

	_, err = DB.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	fmt.Println("Database connected and table created!")
}

// Close closes the database connection
func Close() {
	if DB != nil {
		DB.Close()
	}
}
