package User

import (
	"net/http"

	"github.com/BatuhanSerin/final-project/check"
	"github.com/BatuhanSerin/final-project/internal/api"
	"github.com/BatuhanSerin/final-project/internal/httpErrors"
	"github.com/BatuhanSerin/final-project/package/middleware"
	"github.com/go-openapi/strfmt"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	repo *UserRepository
}

func NewUserHandler(r *gin.RouterGroup, repo *UserRepository, secret string) {
	u := &userHandler{repo: repo}

	r.POST("/login", u.login)
	r.POST("/signup", u.signup)
	r.Use(middleware.Authorization(secret))
	r.POST("/verify", u.VerifyToken)
}

//signup creates a new user
func (u *userHandler) signup(c *gin.Context) {
	var userBody *api.User

	if err := c.Bind(&userBody); err != nil {
		c.JSON(httpErrors.ErrorResponse(httpErrors.CannotBindGivenData))
		return
	}

	if err := userBody.Validate(strfmt.NewFormats()); err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	phoneNumberDigits := check.CheckDigits(userBody.Phone)
	if phoneNumberDigits == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Phone number must be 10 digits",
		})
		return
	}

	user, err := u.repo.createUser(responseToUser(userBody))
	if err != nil {
		c.JSON(httpErrors.ErrorResponse(err))
		return
	}

	c.JSON(http.StatusOK, UserToResponse(user))

	authorizedUser, _ := u.repo.getUser(&user.Email, &user.Password)
	if authorizedUser == nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	}

	_, token, err := u.repo.Login(authorizedUser.Email, authorizedUser.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
	}
	c.JSON(http.StatusOK, token)
}

//login logs in a user
func (u *userHandler) login(c *gin.Context) {
	var user *api.Login

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{
			"message": "Bad Request",
		})
	}

	authorizedUser, _ := u.repo.getUser(user.Email, user.Password)
	if authorizedUser == nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
		})
	}
	userInfo, token, err := u.repo.Login(authorizedUser.Email, authorizedUser.Password)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
		})
	}

	c.JSON(http.StatusOK, userInfo)
	c.JSON(http.StatusOK, token)

}

//VerifyToken verifies the token
func (u *userHandler) VerifyToken(c *gin.Context) {
	u.repo.VerifyToken(c)
	c.JSON(http.StatusOK, gin.H{
		"message": "Token verified",
	})
}
