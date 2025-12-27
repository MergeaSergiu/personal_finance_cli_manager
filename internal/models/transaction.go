package models

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	CategoryID uint
	Amount     float32

	Category Category `gorm:"foreignKey:CategoryID"`
}
