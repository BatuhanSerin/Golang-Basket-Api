package product

import (
	"strconv"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/models"
	page "github.com/BatuhanSerin/final-project/package/pagination"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

//create creates a new product
func (r *ProductRepository) create(product *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.create", zap.Any("product", product))

	if err := r.db.Create(product).Error; err != nil {
		zap.L().Error("product.repo.create Failed", zap.Error(err))
		return nil, err
	}
	return product, nil
}

//createBulk creates new products from uploaded csv file
func (r *ProductRepository) createBulk(csvLines [][]string) ([]models.Product, error) {

	products := []models.Product{}

	for _, line := range csvLines[1:] {

		CategoryId, _ := strconv.ParseInt(line[4], 10, 0)
		Id, _ := strconv.ParseInt(line[0], 10, 0)
		Price, _ := strconv.ParseFloat(line[2], 32)
		Stock, _ := strconv.ParseInt(line[3], 10, 0)

		productBody := &api.Product{
			Category: &api.CategoryWithoutRequired{
				ID: CategoryId,
			},
			ID:    Id,
			Name:  &line[1],
			Price: &Price,
			Stock: &Stock,
		}

		if err := productBody.Validate(strfmt.NewFormats()); err != nil {
			zap.L().Error("product.repo.createBulk Validate Failed", zap.Error(err))
		}

		product, err := r.create(responseToProduct(productBody))
		if err != nil {
			zap.L().Error("product.repo.createBulk Validate Failed", zap.Error(err))
		}
		products = append(products, *product)
	}
	return products, nil

}

//getAll returns all products
func (r *ProductRepository) getAll(c *gin.Context) (*[]models.Product, error) {
	zap.L().Debug("product.repo.getAll")

	var products = &[]models.Product{}

	if err := r.db.Preload("Category").Scopes(page.Paginate(c)).Find(&products).Error; err != nil {
		zap.L().Error("product.repo.getAll Failed", zap.Error(err))
		return nil, err
	}

	return products, nil
}

//getByName returns product by name
func (r *ProductRepository) getByName(name string) (*models.Product, error) {
	zap.L().Debug("product.repo.getByName", zap.Any("name", name))

	var product = &models.Product{}

	if err := r.db.Preload("Category").Where("name LIKE ? ", "%"+name+"%").Find(&product).Error; err != nil {
		zap.L().Error("product.repo.getByName Failed", zap.Error(err))
		return nil, err
	}

	return product, nil
}

//getByID returns product by id
func (r *ProductRepository) getByID(id string) (*models.Product, error) {
	zap.L().Debug("product.repo.getByID", zap.Any("id", id))

	var product = &models.Product{}

	if err := r.db.Preload("Category").First(&product, id).Error; err != nil {
		zap.L().Error("product.repo.getByID Failed", zap.Error(err))
		return nil, err
	}

	return product, nil
}

//update updates product
func (r *ProductRepository) update(product *models.Product) (*models.Product, error) {
	zap.L().Debug("product.repo.update", zap.Any("product", product))

	if err := r.db.Save(&product).Error; err != nil {
		zap.L().Error("product.repo.update Failed", zap.Error(err))
		return nil, err
	}

	return product, nil
}

//delete deletes product
func (r *ProductRepository) delete(product *models.Product) error {
	zap.L().Debug("product.repo.delete", zap.Any("product", product))

	if err := r.db.Delete(&product).Error; err != nil {
		zap.L().Error("product.repo.delete Failed", zap.Error(err))
		return err
	}

	return nil
}

//Migration creates table for products
func (r *ProductRepository) Migration() {
	zap.L().Debug("Product Migration")

	if err := r.db.AutoMigrate(&models.Product{}); err != nil {
		zap.L().Error("Product Migration Failed", zap.Error(err))
	}
}
