package product

import (
	"gorm.io/gorm"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/models"
)

//ProductToResponse converts product to response
func ProductToResponse(p *models.Product) *api.Product {

	return &api.Product{
		Category: &api.CategoryWithoutRequired{
			ID:   int64(p.Category.ID),
			Name: *p.Category.Name,
		},
		ID:    int64(p.ID),
		Name:  p.Name,
		Price: &p.Price,
		Stock: &p.Stock,
	}
}

//ProductToResponseWithoutCategory converts product to response without category
func ProductToResponseWithoutCategory(p *models.Product) *api.Product {

	return &api.Product{
		ID:    int64(p.ID),
		Name:  p.Name,
		Price: &p.Price,
		Stock: &p.Stock,
	}
}

//productsToResponseWithoutCategory converts products to response without category
func productsToResponseWithoutCategory(ps *[]models.Product) []*api.Product {
	products := make([]*api.Product, 0)
	for _, p := range *ps {
		products = append(products, ProductToResponseWithoutCategory(&p))
	}
	return products
}

//productsToResponse converts products to response
func productsToResponse(ps *[]models.Product) []*api.Product {
	products := make([]*api.Product, 0)
	for _, p := range *ps {
		products = append(products, ProductToResponse(&p))
	}
	return products
}

//responseToProduct converts response to product
func responseToProduct(p *api.Product) *models.Product {
	return &models.Product{
		Model:      gorm.Model{ID: uint(p.ID)},
		Name:       p.Name,
		Price:      *p.Price,
		Stock:      *p.Stock,
		CategoryID: uint(p.Category.ID),
	}
}
