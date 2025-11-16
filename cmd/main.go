package main

import (
	"fmt"
	"log"
	"peronal_finance_cli_manager/internal/db"
)

func main() {
	db.Init()
	defer db.Close()

	_, err := db.DB.Exec(`INSERT INTO transactions (date, description, amount, category) VALUES (?, ?, ?, ?)`,
		"2025-11-16", "Coffee", 3.5, "Food")
	if err != nil {
		log.Fatalf("Failed to insert transaction: %v", err)
	}

	// Query and print transactions
	rows, err := db.DB.Query("SELECT id, date, description, amount, category FROM transactions")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	fmt.Println("Transactions:")
	for rows.Next() {
		var id int
		var date, description, category string
		var amount float64
		err = rows.Scan(&id, &date, &description, &amount, &category)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d | %s | %s | %.2f | %s\n", id, date, description, amount, category)
	}
}
