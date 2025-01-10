package routers

import (
	"embed"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetRouters(fiberApp *fiber.App, embedScreens embed.FS) {
	// Routes
	fiberApp.Get("/", func(c *fiber.Ctx) error {
		return controllers.Index(c, embedScreens)
	})

	fiberApp.Get("/hello", controllers.HelloHandler)
	// fiberApp.Get("/input-real-json", controllers.ParsingInputJson)
	// fiberApp.Get("/example", controllers.ExampleHandler)

	// UI GROUP ROUTER
	// ui := fiberApp.Group("/ui")
	// ui.Get("/get-format-validation", controllers.GetFormatValidation)

}
