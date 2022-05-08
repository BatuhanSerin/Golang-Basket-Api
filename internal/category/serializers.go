package category

import (
	"gorm.io/gorm"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/models"

	"github.com/BatuhanSerin/final-project/internal/product"
)

func categoryToResponse(c *models.Category) *api.Category {

	products := make([]*api.Product, 0)
	for _, p := range c.Products {
		products = append(products, product.ProductToResponseWithoutCategory(&p))
	}

	return &api.Category{
		ID:       int64(c.ID),
		Name:     c.Name,
		Products: products,
	}
}

//categoryToResponseWithoutProducts converts category to response without products
func categoryToResponseWithoutProducts(c *models.Category) *api.Category {

	return &api.Category{
		ID:   int64(c.ID),
		Name: c.Name,
	}
}

//categoriesToResponse converts categories to response
func categoriesToResponse(cs *[]models.Category) []*api.Category {
	categories := make([]*api.Category, 0)
	for _, c := range *cs {
		categories = append(categories, categoryToResponse(&c))
	}
	return categories
}

//responseToCategory converts response to category
func responseToCategory(c *api.Category) *models.Category {
	return &models.Category{
		Model: gorm.Model{ID: uint(c.ID)},
		Name:  c.Name,
	}
}
