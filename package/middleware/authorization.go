package middleware

import (
	jwtToken "github.com/BatuhanSerin/final-project/package/jwt"
	"github.com/gin-gonic/gin"
)

func Authorization(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		decodedClaims, _ := jwtToken.ParseToken(token, secret)
		claims := *decodedClaims

		if decodedClaims == nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		//isAdmin
		if claims["role"] != false {
			c.Next()
			c.Abort()
			return
		}
		c.JSON(401, gin.H{
			"message": "You are not an admin",
		})
		c.Abort()
	}
}

func AuthorizationForUser(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		decodedClaims, _ := jwtToken.ParseToken(token, secret)
		claims := *decodedClaims

		if decodedClaims == nil {
			c.JSON(401, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		//isUser
		if claims["userID"] != nil {
			c.Next()
			c.Abort()
			return
		}
		c.JSON(401, gin.H{
			"message": "You are not a User! Please login",
		})
		c.Abort()
	}
}
