package basket

import (
	"fmt"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/models"
	cfg "github.com/BatuhanSerin/final-project/package/config"
	jwtPackage "github.com/BatuhanSerin/final-project/package/jwt"
	"github.com/gin-gonic/gin"
)

//Verify checks if token is valid and set UserID
func Verify(c *gin.Context) *api.Basket {

	tokenString := c.GetHeader("Authorization")
	token, _ := jwtPackage.ParseToken(tokenString, cfg.GetSecretKey())

	claims := *token
	ID := claims["userID"].(float64)

	basketBody := &api.Basket{
		UserID: int64(ID),
	}
	fmt.Printf("%+v\n", basketBody)
	return basketBody

}

//CheckStock checks if there is enough stock for basket
func CheckStock(b *BasketRepository, c *gin.Context, basket *models.Basket) error {
	basketCopy := &models.Basket{
		UserID:    basket.UserID,
		ProductID: basket.ProductID,
		Quantity:  basket.Quantity,
		Products:  basket.Products,
	}
	if err := b.db.Preload("Products").Find(&basketCopy.Products, basketCopy.ProductID).Error; err != nil {
		if uint(basketCopy.Products[0].Stock) < basketCopy.Quantity {
			return fmt.Errorf("Stock is not enough")
		}
	}
	return nil
}

//calculateTotalPrice calculates total price of basket
func calculatePrice(basket *models.Basket) {
	price := float64(basket.Quantity) * basket.Products[0].Price
	basket.TotalPrice = price
}
