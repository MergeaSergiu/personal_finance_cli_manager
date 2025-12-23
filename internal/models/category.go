package models

type Category struct {
	ID     uint    `gorm:"primaryKey"`
	Name   string  `gorm:"unique"`
	Budget float32 `gorm:"not null;default:0"`
}
