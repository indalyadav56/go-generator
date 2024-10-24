package routes

import (
	"backend/internal/todo/controllers"
	"github.com/gin-gonic/gin"
)

func TodoRoutes(router *gin.Engine, c controllers.TodoController) {
	todoV1 := router.Group("/v1/todos")
	{
		todoV1.GET("", c.Get)
		todoV1.POST("", c.Create)
		todoV1.PATCH("", c.Update)
		todoV1.DELETE("", c.Delete)
	}
}
