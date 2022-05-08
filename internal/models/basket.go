package models

import (
	"gorm.io/gorm"
)

type Basket struct {
	gorm.Model
	Quantity   uint
	ProductID  uint
	Products   []Product `gorm:"foreignKey:ID;references:ProductID"`
	TotalPrice float64   `gorm:"type:decimal"`
	UserID     int64     `gorm:"type:bigint"`
}
