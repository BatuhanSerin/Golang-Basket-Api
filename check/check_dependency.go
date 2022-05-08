package check

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

//CheckDependency is a function to check dependency
func CheckDependency(DB *gorm.DB, r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.GET("/ready", func(c *gin.Context) {
		db, err := DB.DB()
		if err != nil {
			zap.L().Fatal("check.dependency.ready Failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err,
			})
			return
		}
		if err := db.Ping(); err != nil {
			zap.L().Fatal("check.dependency.ready Failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  err,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"error":  nil,
		})
	})
}
