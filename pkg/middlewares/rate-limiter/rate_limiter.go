package ratelimiter

import "github.com/gin-gonic/gin"

func RateLimiterMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}
