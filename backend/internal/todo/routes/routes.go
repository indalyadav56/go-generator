package routes

import (
	"backend/internal/todo/controllers"
	"github.com/gofiber/fiber/v2"
)

func TodoRoutes(app *fiber.App, c controllers.TodoController) {
	todoV1 := app.Group("/v1/todos")
	{
		todoV1.Get("", c.Get)
		todoV1.Post("", c.Create)
		todoV1.Patch("", c.Update)
		todoV1.Delete("", c.Delete)
	}
}
