package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"peronal_finance_cli_manager/internal/models"
)

var _ *gorm.DB

func CreateTransaction(categoryName string, amount float32) (*models.Transaction, error) {
	var cat models.Category
	if err := DB.Where("name = ?", categoryName).First(&cat).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category '%s' not found", categoryName)
		}
		return nil, err
	}

	tx := &models.Transaction{
		CategoryID: cat.ID,
		Amount:     amount,
	}

	if err := DB.Create(tx).Error; err != nil {
		return nil, err
	}

	// attach category for convenience
	tx.Category = cat

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
