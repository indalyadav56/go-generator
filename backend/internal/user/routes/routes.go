package routes

import (
	"backend/internal/user/controllers"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, c controllers.UserController) {
	userV1 := app.Group("/v1/users")
	{
		userV1.Get("", c.Get)
		userV1.Post("", c.Create)
		userV1.Patch("", c.Update)
		userV1.Delete("", c.Delete)
	}
}
