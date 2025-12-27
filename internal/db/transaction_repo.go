package db

import (
	"errors"
	"fmt"
	"peronal_finance_cli_manager/internal/models"
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
