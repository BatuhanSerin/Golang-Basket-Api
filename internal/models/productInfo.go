package models

import "gorm.io/gorm"

type ProductInfo struct {
	gorm.Model
	ProductID uint  `gorm:"type:bigint"`
	Quantity  int64 `gorm:"type:bigint"`
	BasketID  uint
	Basket    Basket
	Product   Product `gorm:"foreignkey:ID"`
}

func (ProductInfo) TableName() string { return "product_info" }
