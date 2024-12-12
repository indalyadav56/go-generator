package middlewares

import (
	"common/dto"
	"common/pkg/logger"
	"common/pkg/utils/encryption"
	"common/pkg/utils/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DecryptMiddleware(logger logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("user_id").(string)

		tokenDataParts := strings.Split(userID, "|")
		if len(tokenDataParts) < 2 {
			return
		}

		signingKey := tokenDataParts[1]

		req := new(dto.EncryptedRequest)

		if err := c.ShouldBindJSON(req); err != nil {
			response.HandleErrorWithAbort(c, http.StatusBadRequest, "Invalid request payload", err)
			return
		}

		decrypted, err := encryption.Decrypt(req.Data, []byte(signingKey))
		if err != nil {
			response.HandleErrorWithAbort(c, http.StatusBadRequest, "Invalid request payload data", err)
			return
		}

		c.Set("decrypted", decrypted)

		c.Next()
	}
}
