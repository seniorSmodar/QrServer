package routes

import (
	"workspace/controllers"
	"workspace/middleware"

	"github.com/gofiber/fiber/v2"
)

func VisitRoutes(app *fiber.App) {
	app.Get("/Visit/GainQr", middleware.JWTProtected(), controllers.CreateQr)

	app.Get("/Visits", middleware.JWTProtected(), controllers.GetVisits)

	app.Delete("/Visits/DeleteLegasy", middleware.JWTProtected(), controllers.DeleteLegacyVisits)

	app.Get("/Visits/:code", middleware.JWTProtected(), controllers.CreateVisit)
}