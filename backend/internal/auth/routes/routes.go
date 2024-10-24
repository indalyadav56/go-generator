package routes

import (
	"backend/internal/auth/controllers"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine, c controllers.AuthController) {
	authV1 := router.Group("/v1/auths")
	{
		authV1.GET("", c.Get)
		authV1.POST("", c.Create)
		authV1.PATCH("", c.Update)
		authV1.DELETE("", c.Delete)
	}
}
