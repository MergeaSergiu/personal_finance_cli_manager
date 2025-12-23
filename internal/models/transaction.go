package models

type Transaction struct {
	ID         uint `gorm:"primaryKey"`
	CategoryID uint
	Amount     float64
	Type       string // income or expense
}
