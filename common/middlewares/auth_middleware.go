package middlewares

import (
	"common/pkg/jwt"
	"common/pkg/logger"
	"common/pkg/utils/response"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const BearerTokenPrefix = "Bearer"

func AuthMiddleware(logger logger.Logger, jwt jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractBearerToken(c)
		if err != nil {
			logger.Error(fmt.Sprintf("Authorization error: %s | IP: %s", err.Error(), c.ClientIP()))
			resp := response.Unauthorized(fmt.Sprintf("Authorization error:-%s", err.Error()))
			c.JSON(resp.Status, resp)
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(token)
		if err != nil {
			logger.Error(fmt.Sprintf("Invalid token: %v | IP: %s", err, c.ClientIP()))
			resp := response.Error(http.StatusBadRequest, "Invalid token", err.Error())
			c.JSON(resp.Status, resp)
			c.Abort()
			return
		}

		tokenData, err := jwt.GetClaims(claims)
		if err != nil {
			logger.Error(fmt.Sprintf("Invalid token: %v | IP: %s", err, c.ClientIP()))
			resp := response.Error(http.StatusBadRequest, "Invalid token", err.Error())
			c.JSON(resp.Status, resp)
			c.Abort()
			return
		}

		userId, ok := tokenData["user_id"].(string)
		if !ok {
			resp := response.Error(http.StatusBadRequest, "user id does not available in token", nil)
			c.JSON(resp.Status, resp)
			c.Abort()
			return
		}

		c.Set("user_id", userId)

		c.Next()
	}
}

func extractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != BearerTokenPrefix {
		return "", fmt.Errorf("invalid authorization format")
	}

	return tokenParts[1], nil
}
