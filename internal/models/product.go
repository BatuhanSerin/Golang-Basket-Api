package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID         uint
	Name       *string `gorm:"unique"`
	Price      float64
	Stock      int64
	CategoryID uint
	Category   Category
}
