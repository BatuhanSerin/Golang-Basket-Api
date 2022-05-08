package basket

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BatuhanSerin/final-project/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BasketRepository struct {
	db *gorm.DB
}

func NewBasketRepository(db *gorm.DB) *BasketRepository {
	return &BasketRepository{db: db}
}

func (b *BasketRepository) Migration() {
	zap.L().Debug("basket Migration")

	if err := b.db.AutoMigrate(&models.Basket{}); err != nil {
		zap.L().Error("basket Migration Failed", zap.Error(err))
	}
}

//VerifyToken verifies token
func (b *BasketRepository) VerifyToken(c *gin.Context, basket *models.Basket) (*models.Basket, error) {

	zap.L().Debug("basket.repo.VerifyToken")

	if err := b.db.FirstOrCreate(basket).Error; err != nil {
		zap.L().Error("basket.repo.VerifyToken Failed", zap.Error(err))
		return nil, err
	}

	return basket, nil

}

//Increment increments quantity of product
func (b *BasketRepository) Increment(c *gin.Context, basket *models.Basket) (*models.Basket, error) {

	zap.L().Debug("basket.repo.Increment")

	if err := b.db.Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Where(fmt.Sprintf("product_id = %v", basket.ProductID)).First(basket).Error; err != nil {
		basket.Quantity = 1
		if err := b.db.Create(basket).Error; err != nil {
			zap.L().Error("basket.repo.Increment Failed", zap.Error(err))
			return nil, err
		}
	} else {

		basket.Quantity++
		quantity := strconv.FormatUint(uint64(basket.Quantity), 10)
		//CheckStock checks stock of product is enough for quantity
		if err := CheckStock(b, c, basket); err != nil {
			return nil, err
		}

		if err := b.db.Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Where(fmt.Sprintf("product_id = %v", basket.ProductID)).Update("quantity", quantity).Error; err != nil {
			zap.L().Error("basket.repo.Increment Failed", zap.Error(err))
			return nil, err
		}
	}
	return basket, nil

}

//Decrement decrements quantity of product
func (b *BasketRepository) Decrement(c *gin.Context, basket *models.Basket) (*models.Basket, error) {

	zap.L().Debug("basket.repo.Decrement")

	if err := b.db.Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Where(fmt.Sprintf("product_id = %v", basket.ProductID)).First(basket).Error; err == nil {
		basket.Quantity--
		if basket.Quantity > 0 {
			quantity := strconv.FormatUint(uint64(basket.Quantity), 10)
			if err := b.db.Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Where(fmt.Sprintf("product_id = %v", basket.ProductID)).Update("quantity", quantity).Error; err != nil {
				zap.L().Error("basket.repo.Decrement Failed", zap.Error(err))
				return nil, err
			}
		} else {
			if err := b.db.Delete(basket).Error; err != nil {
				zap.L().Error("basket.repo.Decrement Failed", zap.Error(err))
				return nil, err
			}
			return nil, err
		}
	}

	return basket, nil

}

//Create creates new basket
func (b *BasketRepository) Create(c *gin.Context, basket *models.Basket) (*models.Basket, error) {

	//CheckStock checks stock of product is enough for quantity
	if err := CheckStock(b, c, basket); err != nil {
		return nil, err
	}

	zap.L().Debug("basket.repo.Create")
	quantity := strconv.FormatUint(uint64(basket.Quantity), 10)
	if err := b.db.Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Where(fmt.Sprintf("product_id = %v", basket.ProductID)).Update("quantity", quantity).First(basket).Error; err != nil {
		if err := b.db.Create(basket).Error; err != nil {
			zap.L().Error("basket.repo.Create Failed", zap.Error(err))
			return nil, err
		}
	}
	return basket, nil
}

//ListCartItems returns list of cart items
func (b *BasketRepository) ListCartItems(c *gin.Context, basket *models.Basket) ([]*models.Basket, error) {

	zap.L().Debug("basket.repo.ListCartItems")

	rows, _ := b.db.Preload("Products").Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Rows()
	defer rows.Close()

	items := make([]*models.Basket, 0)

	for rows.Next() {
		var basketCopy *models.Basket
		// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		b.db.ScanRows(rows, &basketCopy)

		if err := b.db.Preload("Products").First(basketCopy, basketCopy.ID).Error; err != nil {
			zap.L().Error("basket.repo.ListCartItems Failed", zap.Error(err))
			return nil, err
		}
		calculatePrice(basketCopy)
		items = append(items, basketCopy)
	}

	return items, nil
}

//Buy buys items
func (b *BasketRepository) Buy(c *gin.Context, basket *models.Basket) ([]*models.Basket, error) {

	zap.L().Debug("basket.repo.Buy")

	items := make([]*models.Basket, 0)

	//Checks if basket has at least one item

	if err := b.db.Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Where(fmt.Sprintf("product_id = %v", basket.ProductID)).First(basket).Error; err == nil {
		//If basket is empty
		return nil, nil
	} else {
		rows, _ := b.db.Preload("Products").Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Rows()
		defer rows.Close()

		for rows.Next() {
			var basketCopy *models.Basket

			b.db.ScanRows(rows, &basketCopy)

			if err := b.db.Preload("Products").First(basketCopy, basketCopy.ID).Error; err != nil {
				zap.L().Error("basket.repo.ListCartItems Failed", zap.Error(err))
				return nil, err
			} else {
				product := basketCopy.Products[0]
				stock := product.Stock - int64(basketCopy.Quantity)
				if err := b.db.Model(&product).Update("Stock", stock).Error; err != nil {
					zap.L().Error("basket.repo.ListCartItems Failed", zap.Error(err))
					return nil, err
				}
				if err := b.db.Delete(basketCopy).Error; err != nil {
					zap.L().Error("basket.repo.ListCartItems Failed", zap.Error(err))
					return nil, err
				}
			}
			calculatePrice(basketCopy)
			items = append(items, basketCopy)
		}

		return items, nil
	}
}

//Order orders items
func (b *BasketRepository) Order(c *gin.Context, basket *models.Basket) ([]*models.Basket, error) {

	zap.L().Debug("basket.repo.Order")

	rows, _ := b.db.Preload("Products").Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Unscoped().Rows()
	defer rows.Close()

	items := make([]*models.Basket, 0)

	for rows.Next() {
		var basketCopy *models.Basket
		// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		b.db.ScanRows(rows, &basketCopy)

		if err := b.db.Preload("Products").Unscoped().First(basketCopy, basketCopy.ID).Error; err != nil {
			zap.L().Error("basket.repo.ListCartItems Failed", zap.Error(err))
			return nil, err
		}
		if basketCopy.DeletedAt.Valid == true {

			calculatePrice(basketCopy)
			items = append(items, basketCopy)
		}
	}
	return items, nil

}

//Cancel cancels items
func (b *BasketRepository) Cancel(c *gin.Context, basket *models.Basket) ([]*models.Basket, error) {

	zap.L().Debug("basket.repo.Cancel")
	t := time.Now().UTC()
	fmt.Printf("\n\n\n14ten buyukse\n%v\n\n\n", gorm.DeletedAt{})

	rows, _ := b.db.Preload("Products").Model(&basket).Where(fmt.Sprintf("user_id = %v", basket.UserID)).Unscoped().Rows()
	defer rows.Close()

	items := make([]*models.Basket, 0)

	for rows.Next() {
		var basketCopy *models.Basket
		// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		b.db.ScanRows(rows, &basketCopy)

		if err := b.db.Preload("Products").Unscoped().First(basketCopy, basketCopy.ID).Error; err != nil {
			zap.L().Error("basket.repo.ListCartItems Failed", zap.Error(err))
			return nil, err
		}
		if basketCopy.DeletedAt.Valid == true {

			diff := t.Sub(basketCopy.DeletedAt.Time)
			hour := int(diff.Hours()) * -1
			fmt.Printf("\n\n\n14ten buyukse\n%v\n\n\n", hour)
			if hour < 336 {

				newbasket := models.Basket{
					UserID:    basketCopy.UserID,
					ProductID: basketCopy.ProductID,
					Quantity:  basketCopy.Quantity,
				}

				if err := b.db.Create(&newbasket).Error; err != nil {
					zap.L().Error("basket.repo.Create Failed", zap.Error(err))
					return nil, err
				}
				calculatePrice(basketCopy)
				items = append(items, basketCopy)
			}
		}
	}
	return items, nil

}

func (b *BasketRepository) GetByID(c *gin.Context, id string) (*models.Basket, error) {

	zap.L().Debug("basket.repo.GetByID")

	basket := &models.Basket{}

	if err := b.db.Preload("Products").First(basket, id).Error; err != nil {
		zap.L().Error("basket.repo.GetByID Failed", zap.Error(err))
		return nil, err
	}

	return basket, nil

}
