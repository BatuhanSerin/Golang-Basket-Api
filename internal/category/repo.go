package category

import (
	"github.com/BatuhanSerin/final-project/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

//create creates a new category
func (r *CategoryRepository) create(category *models.Category) (*models.Category, error) {
	zap.L().Debug("category.repo.create", zap.Any("category", category))

	if err := r.db.Create(category).Error; err != nil {
		zap.L().Error("category.repo.create Failed", zap.Error(err))
		return nil, err
	}
	return category, nil
}

//getByID gets a category by id
func (r *CategoryRepository) getByID(id string) (*models.Category, error) {
	zap.L().Debug("Get category By ID", zap.Any("id", id))

	var category = &models.Category{}

	if err := r.db.Preload("Products").First(&category, id).Error; err != nil {
		zap.L().Error("Get category By ID Failed", zap.Error(err))
		return nil, err
	}

	return category, nil
}

//update updates a category
func (r *CategoryRepository) update(category *models.Category) (*models.Category, error) {
	zap.L().Debug("Update category", zap.Any("category", category))

	if err := r.db.Save(&category).Error; err != nil {
		zap.L().Error("Update category Failed", zap.Error(err))
		return nil, err
	}

	return category, nil
}

//delete deletes a category
func (r *CategoryRepository) delete(id string) error {
	zap.L().Debug("Delete category", zap.Any("id", id))

	category, err := r.getByID(id)

	if err != nil {
		return err
	}

	if err := r.db.Delete(&category).Error; err != nil {
		zap.L().Error("Delete category Failed", zap.Error(err))
		return err
	}

	return nil
}

//Migrate migrates the database
func (r *CategoryRepository) Migration() {
	zap.L().Debug("category Migration")

	if err := r.db.AutoMigrate(&models.Category{}); err != nil {
		zap.L().Error("category Migration Failed", zap.Error(err))
	}
}
