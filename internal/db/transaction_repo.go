package db

import (
	"errors"
	"fmt"
	"peronal_finance_cli_manager/internal/models"
	"peronal_finance_cli_manager/internal/transaction"
	"time"

	"gorm.io/gorm"
)

var _ *gorm.DB

func CreateTransaction(categoryName string, amount float32, dateStr string) (*models.Transaction, error) {
	var cat models.Category
	if err := DB.Where("name = ?", categoryName).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category '%s' not found", categoryName)
		}
		return nil, err
	}

	// parse date string
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}

	tx := &models.Transaction{
		CategoryID: cat.ID,
		Amount:     amount,
		Date:       date,
	}

	if err := DB.Create(tx).Error; err != nil {
		return nil, err
	}

	// attach category for convenience
	tx.Category = cat

	// Check budget
	if err := CheckBudget(DB, cat, amount, dateStr); err != nil {
		fmt.Println("Budget alert triggered")
	}

	return tx, nil
}

func GetTransactionsByCategory(categoryID uint) ([]models.Transaction, error) {
	var txs []models.Transaction

	err := DB.
		Preload("Category").
		Where("category_id = ?", categoryID).
		Order("id DESC").
		Find(&txs).Error

	return txs, err
}

// ImportTransactionsFromFile parses a file (CSV/OFX) and inserts transactions into the DB.
func ImportTransactionsFromFile(filePath string) ([]models.Transaction, error) {
	var transactions []models.Transaction
	var err error

	switch transaction.DetectFormat(filePath) {
	case "csv":
		transactions, err = transaction.ParseCSV(filePath)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported file format")
	}

	var imported []models.Transaction
	for _, tx := range transactions {
		newTx, err := CreateTransaction(tx.Category.Name, tx.Amount, tx.Date.Format("2006-01-02"))
		if err != nil {
			// Skip invalid transactions but log error
			fmt.Printf("Failed to import transaction: %v\n", err)
			continue
		}
		imported = append(imported, *newTx)
	}

	return imported, nil
}
