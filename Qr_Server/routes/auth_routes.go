package routes

import (
	"workspace/controllers"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app *fiber.App) {
	app.Post("/Auth/Token", controllers.CreateToken)
}