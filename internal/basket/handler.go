package basket

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/httpErrors"
	cfg "github.com/BatuhanSerin/final-project/package/config"
	jwtPackage "github.com/BatuhanSerin/final-project/package/jwt"
	"github.com/BatuhanSerin/final-project/package/middleware"

	"github.com/gin-gonic/gin"
)

type basketHandler struct {
	repo *BasketRepository
}

func NewBasketHandler(r *gin.RouterGroup, repo *BasketRepository, secret string) {
	b := &basketHandler{repo: repo}

	r.Use(middleware.AuthorizationForUser(secret))
	r.POST("/verify", b.VerifyToken)
	r.POST("/create/:id/:quantity", b.create)
	r.POST("/inc/:id", b.increment)
	r.POST("/dec/:id", b.decrement)
	r.GET("/list", b.listCartItems)
	r.GET("/buy", b.buy)
	r.GET("/order", b.order)
	r.GET("cancel", b.cancel)
	r.GET("/:id", b.getByID)
}

//VerifyToken is a middleware to verify token
func (b *basketHandler) VerifyToken(c *gin.Context) {

	tokenString := c.GetHeader("Authorization")
	token, err := jwtPackage.ParseToken(tokenString, cfg.GetSecretKey())
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Authorized",
	})

	c.JSON(http.StatusOK, token)

	claims := *token
	ID := claims["userID"].(float64)

	basketBody := &api.Basket{
		UserID: int64(ID),
	}

	basket, err := b.repo.VerifyToken(c, responseToBasket(basketBody))
	c.JSON(http.StatusOK, gin.H{
		"message": "Token verified",
	})
	c.JSON(http.StatusOK, basketToResponse(basket))
}

//create creates a new basket
func (b *basketHandler) create(c *gin.Context) {

	basketBody := Verify(c)
	mux := &sync.RWMutex{}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 0)
	quantity, _ := strconv.ParseInt(c.Param("quantity"), 10, 0)

	basketBody.ProductID = int64(id)
	basketBody.Quantity = int64(quantity)

	mux.Lock()

	basket, err := b.repo.Create(c, responseToBasket(basketBody))

	mux.Unlock()

	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		c.JSON(http.StatusOK, gin.H{
			"message": "Stock is NOT enough",
		})
		return
	}
	c.JSON(http.StatusOK, basketToResponse(basket))
}

//increment increments a product quantity by 1
func (b *basketHandler) increment(c *gin.Context) {
	basketBody := Verify(c)
	mux := &sync.RWMutex{}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 0)
	basketBody.ProductID = int64(id)

	mux.Lock()

	basket, err := b.repo.Increment(c, responseToBasket(basketBody))

	mux.Unlock()

	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		c.JSON(http.StatusOK, gin.H{
			"message": "Stock is NOT enough",
		})
		return
	}
	c.JSON(http.StatusOK, basketToResponse(basket))
}

//decrement decrements a product quantity by 1
func (b *basketHandler) decrement(c *gin.Context) {
	basketBody := Verify(c)

	mux := &sync.RWMutex{}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 0)
	basketBody.ProductID = int64(id)

	mux.Lock()

	basket, err := b.repo.Decrement(c, responseToBasket(basketBody))

	mux.Unlock()

	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, basketToResponse(basket))
}

//listCartItems lists all items on user basket
func (b *basketHandler) listCartItems(c *gin.Context) {
	basketBody := Verify(c)
	totalPrice := []float64{0} //slice of float64 to use out of scope

	basket, err := b.repo.ListCartItems(c, responseToBasket(basketBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	for _, v := range basket {
		c.JSON(http.StatusOK, basketToResponse(v))
		totalPrice[0] = totalPrice[0] + v.TotalPrice
	}

	c.JSON(http.StatusOK, gin.H{
		"The Total Price is": fmt.Sprintf("%v", totalPrice[0]),
	})

}

//buy buys all items on user basket
func (b *basketHandler) buy(c *gin.Context) {
	basketBody := Verify(c)
	totalPrice := []float64{0} //slice of float64 to use out of scope
	mux := &sync.RWMutex{}
	mux.Lock()

	basket, err := b.repo.Buy(c, responseToBasket(basketBody))

	mux.Unlock()

	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	for _, v := range basket {
		c.JSON(http.StatusOK, basketToResponse(v))
		totalPrice[0] = totalPrice[0] + v.TotalPrice
	}

	c.JSON(http.StatusOK, gin.H{
		"The Total Price is": fmt.Sprintf("%v", totalPrice[0]),
	})

}

//order orders all items on user basket
func (b *basketHandler) order(c *gin.Context) {
	basketBody := Verify(c)
	totalPrice := []float64{0} //slice of float64 to use out of scope
	mux := &sync.RWMutex{}
	mux.Lock()

	basket, err := b.repo.Order(c, responseToBasket(basketBody))

	mux.Unlock()

	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	for _, v := range basket {
		c.JSON(http.StatusOK, basketToResponse(v))
		totalPrice[0] = totalPrice[0] + v.TotalPrice
	}

	c.JSON(http.StatusOK, gin.H{
		"The Total Price is": fmt.Sprintf("%v", totalPrice[0]),
	})

}

//cancel cancels order if 14 days have been not past
func (b *basketHandler) cancel(c *gin.Context) {
	basketBody := Verify(c)
	totalPrice := []float64{0} //slice of float64 to use out of scope
	mux := &sync.RWMutex{}
	mux.Lock()

	basket, err := b.repo.Cancel(c, responseToBasket(basketBody))

	mux.Unlock()
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Canceled order(s)": "\n",
	})

	for _, v := range basket {
		c.JSON(http.StatusOK, basketToResponse(v))
	}
	if len(basket) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"The Total Price is": fmt.Sprintf("%v", totalPrice[0]),
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "You Can Not Cancel An Order That Has Ordered before 14 days",
		})
	}

}

//getByID gets a basket by id
func (b *basketHandler) getByID(c *gin.Context) {
	id := c.Param("id")
	basket, err := b.repo.GetByID(c, id)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, basketToResponse(basket))
}
