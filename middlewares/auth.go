package middlewares

import (
	"github.com/gin-gonic/gin"
	"library_app/auth"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Request need access token",
				"status":  "Unauthorized",
			})

			c.Abort()
			return
		}
		err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
				"status":  "Unauthorized",
			})

			c.Abort()
			return
		}
		c.Next()
	}
}