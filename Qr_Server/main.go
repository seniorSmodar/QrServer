package main

import (
	"workspace/configs"
	"workspace/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data":"Hello from QrServer"})
	})

	routes.UserRoutes(app)

	routes.AuthRouter(app)

	routes.VisitRoutes(app)

	app.Listen(configs.EnvPort())
}