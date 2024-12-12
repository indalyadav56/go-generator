package middlewares

import (
	"common/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ServiceTokenMiddleware(jwt jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("X-Service-Token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing service token"})
			c.Abort()
			return
		}

		// Parse and validate JWT token
		token, err := jwt.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid service token"})
			c.Abort()
			return
		}

		// // Extract service ID from claims
		// if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 	c.Set("service_id", claims["service_id"])
		// }

		c.Next()
	}
}
