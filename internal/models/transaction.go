package models

import "time"

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	CategoryID uint
	Amount     float32
	Date       time.Time `gorm:"type:date"`

	Category Category `gorm:"foreignKey:CategoryID"`
}
