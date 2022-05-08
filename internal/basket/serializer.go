package basket

import (
	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/models"
	"github.com/BatuhanSerin/final-project/internal/product"
)

//basketToResponse converts basket to response
func basketToResponse(b *models.Basket) *api.Basket {

	products := make([]*api.Product, 0)

	for _, p := range b.Products {
		products = append(products, product.ProductToResponseWithoutCategory(&p))
	}

	return &api.Basket{
		UserID:     b.UserID,
		ProductID:  int64(b.ProductID),
		Quantity:   int64(b.Quantity),
		TotalPrice: b.TotalPrice,
		Products:   products,
	}
}

//responseToBasket converts response to basket
func responseToBasket(b *api.Basket) *models.Basket {

	return &models.Basket{
		Quantity:  uint(b.Quantity),
		UserID:    b.UserID,
		ProductID: uint(b.ProductID),
	}
}
