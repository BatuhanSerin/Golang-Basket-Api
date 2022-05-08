package product

import (
	"encoding/csv"
	"net/http"
	"sync"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/httpErrors"
	"github.com/BatuhanSerin/final-project/package/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type productHandler struct {
	repo *ProductRepository
}

func NewProductHandler(r *gin.RouterGroup, repo *ProductRepository, secret string) {
	p := &productHandler{repo: repo}
	r.GET("/", p.getAll)
	r.GET("/:id", p.getByID)
	r.GET("/search/:name", p.getByName)
	r.Use(middleware.Authorization(secret))
	r.POST("/create", p.create)
	r.POST("/createBulk", p.createBulk)

}

//getByName is a function to get product by name
func (p *productHandler) getByName(c *gin.Context) {
	name := c.Param("name")

	product, err := p.repo.getByName(name)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ProductToResponseWithoutCategory(product))
}

//createBulk creates new products from uploaded csv file
func (p *productHandler) createBulk(c *gin.Context) {

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
	}

	csvFile, err := file.Open()
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
	}

	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
	}
	mux := &sync.RWMutex{}
	mux.Lock()

	products, err := p.repo.createBulk(csvLines)

	mux.Unlock()
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
	}
	c.JSON(http.StatusOK, productsToResponseWithoutCategory(&products))

}

//create creates a new product
func (p *productHandler) create(c *gin.Context) {
	productBody := &api.Product{}

	if err := c.Bind(&productBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := productBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	mux := &sync.RWMutex{}
	mux.Lock()

	product, err := p.repo.create(responseToProduct(productBody))

	mux.Unlock()
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, ProductToResponseWithoutCategory(product))

}

//getAll is a function to get all products
func (p *productHandler) getAll(c *gin.Context) {
	products, err := p.repo.getAll(c)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, productsToResponseWithoutCategory(products))
}

//getByID is a function to get product by id
func (p *productHandler) getByID(c *gin.Context) {
	id := c.Param("id")
	product, err := p.repo.getByID(id)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, ProductToResponseWithoutCategory(product))
}
