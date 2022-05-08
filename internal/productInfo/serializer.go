package productInfo

import (
	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/models"
)

func productInfoToResponse(p *models.ProductInfo) *api.ProductInfo {
	id := int64(p.ProductID)
	return &api.ProductInfo{
		Product: &api.Product{
			ID:    int64(p.Product.ID),
			Stock: &p.Product.Stock,
			Price: &p.Product.Price,
		},
		Basket: &api.BasketWithoutRequired{
			UserID:     int64(p.BasketID),
			TotalPrice: p.Basket.TotalPrice,
		},
		ProductsIDs: &id,
		Quantity:    p.Quantity,
	}
}

func responseToProductInfo(p *api.ProductInfo) *models.ProductInfo {
	return &models.ProductInfo{
		ProductID: uint(*p.ProductsIDs),
		Quantity:  p.Quantity,
		BasketID:  uint(p.Basket.UserID),
		Product: models.Product{
			ID: uint(*p.ProductsIDs),
		},
	}
}
