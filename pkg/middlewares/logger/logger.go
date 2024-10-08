package logger

import "github.com/gin-gonic/gin"

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
