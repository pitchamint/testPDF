package routes

import (
	"github.com/divrhino/fruitful-pdf/controllers"
	"github.com/gofiber/fiber/v2"
)

func Route(app *fiber.App) {
	app.Post("/", controllers.CreateCertificate)
	app.Post("/csv", controllers.CreateCertificateCSV)
}
