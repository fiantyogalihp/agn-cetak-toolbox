package routers

import (
	"embed"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/controllers"
	"github.com/gofiber/fiber/v2"
)

func SetRouters(fiberApp *fiber.App, embedScreens embed.FS) {
	// Routes
	fiberApp.Get("/", controllers.Index)

	// UI COMPONENTS GROUP ROUTER
	components := fiberApp.Group("/components")
	components.Get("/screen-choice", func(c *fiber.Ctx) error {
		return controllers.GetScreenChoices(c, embedScreens)
	})

	// INPUT VALIDATION GROUP ROUTER
	validate := fiberApp.Group("/validate")
	validate.Post("/json-field", func(c *fiber.Ctx) error {
		return controllers.ValidateJSONField(c, embedScreens)
	})

}
