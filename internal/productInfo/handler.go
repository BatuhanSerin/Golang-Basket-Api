package productInfo

import (
	"net/http"

	"github.com/BatuhanSerin/final-project/internal/httpErrors"
	"github.com/BatuhanSerin/final-project/package/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type productInfoHandler struct {
	repo *ProductInfoRepository
}

func NewProductInfoHandler(r *gin.RouterGroup, repo *ProductInfoRepository, secret string) {
	p := &productInfoHandler{repo: repo}
	r.Use(middleware.AuthorizationForUser(secret))
	r.POST("/update", p.update) //Update product info with given quantity
	r.POST("/create", p.create) //Create product if it doesn't exist
	r.POST("/add", p.add)
	r.POST("/dec", p.dec)
	r.GET("/get", p.get)
	// r.GET("/:id", p.getByID)
	// r.GET("/search/:name", p.getByName)
	// r.Use(middleware.Authorization(secret))
	// r.POST("/create", p.create)
	// r.POST("/createBulk", p.createBulk)
}
func (p *productInfoHandler) get(c *gin.Context) {
	category, err := p.repo.get("3")
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, productInfoToResponse(category))
}
func (p *productInfoHandler) dec(c *gin.Context) {

	productInfoBody := initUserId(c)

	if err := c.Bind(&productInfoBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := productInfoBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, productInfoBody)
	c.JSON(http.StatusOK, responseToProductInfo(productInfoBody))
	productInfo, err := p.repo.dec(responseToProductInfo(productInfoBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, productInfoToResponse(productInfo))

}

func (p *productInfoHandler) add(c *gin.Context) {

	productInfoBody := initUserId(c)

	if err := c.Bind(&productInfoBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := productInfoBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	productInfo, err := p.repo.add(responseToProductInfo(productInfoBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, productInfoToResponse(productInfo))

}

func (p *productInfoHandler) create(c *gin.Context) {

	productInfoBody := initUserId(c)

	if err := c.Bind(&productInfoBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := productInfoBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	productInfo, err := p.repo.create(responseToProductInfo(productInfoBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, productInfoToResponse(productInfo))

}
func (p *productInfoHandler) update(c *gin.Context) {

	productInfoBody := initUserId(c)

	if err := c.Bind(&productInfoBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := productInfoBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	productInfo, err := p.repo.update(responseToProductInfo(productInfoBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, productInfoToResponse(productInfo))

}
