package db

import (
	"errors"
	_ "fmt"
	"peronal_finance_cli_manager/internal/models"
)

func CreateCategory(name string) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("Category name is empty")
	}

	cat := models.Category{
		Name: name,
	}
	if err := DB.Create(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func GetCategory(id uint) (*models.Category, error) {
	var cat models.Category
	if err := DB.Where("id = ?", id).First(&cat).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	if err := DB.Find(&categories).Order("name").Error; err != nil {
		return nil, err
	}

	return categories, nil
}
