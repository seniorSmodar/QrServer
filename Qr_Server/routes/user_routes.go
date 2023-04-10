package routes

import (
	"workspace/controllers"
	"workspace/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	app.Post("/User/Register", controllers.Register)

	app.Get("/Users", middleware.JWTProtected(), controllers.GetUsers)

	app.Delete("/User/Delete/:userId", middleware.JWTProtected(), controllers.DeleteUser)
}