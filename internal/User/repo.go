package User

import (
	"net/http"
	"time"

	"github.com/BatuhanSerin/final-project/internal/models"
	cfg "github.com/BatuhanSerin/final-project/package/config"
	jwtPackage "github.com/BatuhanSerin/final-project/package/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

//createUser is a function to create a new user
func (u *UserRepository) createUser(user *models.User) (*models.User, error) {
	zap.L().Debug("user.repo.create", zap.Any("user", user))

	if err := u.db.Create(user).Error; err != nil {
		zap.L().Error("user.repo.create Failed", zap.Error(err))
		return nil, err
	}
	return user, nil
}

//getUser gets a user by email and password
func (u *UserRepository) getUser(email, password *string) (*models.User, error) {
	users := []*models.User{}
	pass := *password
	emailValue := *email

	if err := u.db.Find(&users).Error; err != nil {
		zap.L().Error("user.repo.getUsers Failed", zap.Error(err))
		return nil, err
	}
	zap.L().Debug("user.repo.getUser", zap.Any("email", email), zap.Any("password", password))

	for _, v := range users {

		if v.Email == emailValue && v.Password == pass {

			return v, nil
		}
	}
	return nil, nil
}

//Login is a function to login a user
func (r *UserRepository) Login(email, password string) (jwt.Claims, string, error) {
	var user *models.User
	if err := r.db.Where("email = ? AND password = ?", email, password).First(&user).Error; err != nil {
		return nil, "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"email":  user.Email,
		"role":   user.IsAdmin,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := jwtPackage.GenerateToken(token, cfg.GetSecretKey())

	return token.Claims, tokenString, nil
}

//VerifyToken is a function to verify token
func (r *UserRepository) VerifyToken(c *gin.Context) {

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

	return

}

//Migration is a function to migrate the database
func (u *UserRepository) Migration() {
	zap.L().Debug("user Migration")

	if err := u.db.AutoMigrate(&models.User{}); err != nil {
		zap.L().Error("user Migration Failed", zap.Error(err))
	}
}
