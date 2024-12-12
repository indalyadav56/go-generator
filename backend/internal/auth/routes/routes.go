package routes

import (
	"backend/internal/auth/controllers"
	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App, c controllers.AuthController) {
	authV1 := app.Group("/v1/auths")
	{
		authV1.Get("", c.Get)
		authV1.Post("", c.Create)
		authV1.Patch("", c.Update)
		authV1.Delete("", c.Delete)
	}
}
