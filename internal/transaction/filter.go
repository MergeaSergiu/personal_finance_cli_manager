package transaction

import (
	"peronal_finance_cli_manager/internal/models"
	"time"
)

// FilterByExactDate returns transactions on the exact date
func FilterByExactDate(txs []models.Transaction, date time.Time) []models.Transaction {
	var filtered []models.Transaction
	for _, tx := range txs {
		if tx.Date.Format("2006-01-02") == date.Format("2006-01-02") {
			filtered = append(filtered, tx)
		}
	}
	return filtered
}

// FilterBeforeDate returns transactions before a given date
func FilterBeforeDate(txs []models.Transaction, date time.Time) []models.Transaction {
	var filtered []models.Transaction
	for _, tx := range txs {
		if tx.Date.Before(date) {
			filtered = append(filtered, tx)
		}
	}
	return filtered
}

// FilterByYear returns transactions that occurred in a specific year
func FilterByYear(txs []models.Transaction, year int) []models.Transaction {
	var filtered []models.Transaction
	for _, tx := range txs {
		if tx.Date.Year() == year {
			filtered = append(filtered, tx)
		}
	}
	return filtered
}
