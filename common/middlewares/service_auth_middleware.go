package middlewares

import (
	"common/pkg/jwt"
	"common/pkg/logger"
	"common/pkg/utils/response"
	"crypto/subtle"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ServiceAuthHeader = "X-Service-API-Key"
)

func ServiceAuthMiddleware(logger logger.Logger, jwt jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader(ServiceAuthHeader)
		if apiKey == "" {
			response.HandleUnAuthorizedErrorWithAbort(c, "missing service authentication header")
			return
		}

		// expectedAPIKey := os.Getenv("SERVICE_API_KEY")
		expectedAPIKey := "SERVICE_API_KEY"
		if subtle.ConstantTimeCompare([]byte(apiKey), []byte(expectedAPIKey)) != 1 {
			response.HandleErrorWithAbort(c, http.StatusBadRequest, "invalid service authentication token")
			return
		}

		// Generate JWT token for subsequent requests
		token, err := jwt.GenerateToken(map[string]interface{}{
			"service": "authorized-service",
			"exp":     time.Now().Add(time.Hour * 1).Unix(),
		})
		if err != nil {
			response.HandleServerErrorWithAbort(c, "failed to generate token")
			return
		}

		c.Header("X-Service-Token", token)
		c.Next()
	}
}
