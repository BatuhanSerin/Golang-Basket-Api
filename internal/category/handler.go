package category

import (
	"net/http"
	"strconv"

	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/httpErrors"
	"github.com/BatuhanSerin/final-project/package/middleware"
	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
)

type categoryHandler struct {
	repo *CategoryRepository
}

func NewCategoryHandler(r *gin.RouterGroup, repo *CategoryRepository, secret string) {
	c := &categoryHandler{repo: repo}

	//r.GET("/", h.getAll)
	r.GET("/:id", c.getByID)
	r.Use(middleware.Authorization(secret))
	r.POST("/create", c.create)
	r.PUT("/:id", c.update)
	r.DELETE("/:id", c.delete)
}

//create creates a new category
func (ct *categoryHandler) create(c *gin.Context) {
	categoryBody := &api.Category{}

	if err := c.Bind(&categoryBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := categoryBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	category, err := ct.repo.create(responseToCategory(categoryBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, categoryToResponse(category))

}

//getByID gets a category by id
func (ct *categoryHandler) getByID(c *gin.Context) {
	id := c.Param("id")
	category, err := ct.repo.getByID(id)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, categoryToResponse(category))
}

//update updates a category
func (ct *categoryHandler) update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	categoryBody := &api.Category{ID: int64(id)}

	if err := c.Bind(&categoryBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := categoryBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	category, err := ct.repo.update(responseToCategory(categoryBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, categoryToResponse(category))
}

//delete deletes a category
func (ct *categoryHandler) delete(c *gin.Context) {
	id := c.Param("id")
	category, err := ct.repo.getByID(id)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}
	c.JSON(http.StatusOK, "Deleted Category")
	c.JSON(http.StatusOK, categoryToResponse(category))

	err = ct.repo.delete(id)
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

}
