package main

import (
	"github.com/divrhino/fruitful-pdf/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	routes.Route(app)
	app.Listen(":3000")
}
