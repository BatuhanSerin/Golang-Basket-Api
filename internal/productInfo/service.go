package productInfo

import (
	"net/http"

	"github.com/BatuhanSerin/final-project/internal/api"
	cfg "github.com/BatuhanSerin/final-project/package/config"
	jwtPackage "github.com/BatuhanSerin/final-project/package/jwt"
	"github.com/gin-gonic/gin"
)

// func check(c *gin.Context) error {
// 	productInfoBody := api.ProductInfo{}
// 	if err := c.Bind(&productInfoBody); err != nil {
// 		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
// 		return nil
// 	}

// 	if err := productInfoBody.Validate(strfmt.NewFormats()); err != nil {
// 		c.JSON(httpErrors.ErrorResponse(err))
// 		return nil
// 	}
// }

func initUserId(c *gin.Context) *api.ProductInfo {
	tokenString := c.GetHeader("Authorization")
	token, err := jwtPackage.ParseToken(tokenString, cfg.GetSecretKey())
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return nil
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Authorized",
	})

	c.JSON(http.StatusOK, token)

	claims := *token
	ID := claims["userID"].(float64)

	productInfoBody := &api.ProductInfo{
		Basket: &api.BasketWithoutRequired{UserID: int64(ID)},
	}
	return productInfoBody
}
