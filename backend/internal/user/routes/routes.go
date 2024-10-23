package routes

import (
	"backend/internal/user/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine, c controllers.UserController) {
	userV1 := router.Group("/v1/users")
	{
		userV1.GET("", c.Get)
		userV1.POST("", c.Create)
		userV1.PATCH("", c.Update)
		userV1.DELETE("", c.Delete)
	}
}
