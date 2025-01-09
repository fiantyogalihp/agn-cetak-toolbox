package routers

import (
	"github.com/fiantyogalihp/dynamic-json-parsing-struct/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetRouters(fiberApp *fiber.App) {
	// Routes
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return c.Render("templates/index", fiber.Map{
			"Title": "HTMX + Fiber Quickstart",
		})
	})

	fiberApp.Get("/hello", controllers.HelloHandler)
	// fiberApp.Get("/input-real-json", controllers.ParsingInputJson)
	// fiberApp.Get("/example", controllers.ExampleHandler)

}
