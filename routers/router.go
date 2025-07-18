package routers

import (
	"embed"

	"github.com/fiantyogalihp/agn-cetak-toolbox/controllers"
	"github.com/gofiber/fiber/v2"
)

// Routes
func SetRouters(fiberApp *fiber.App, embedScreens embed.FS) {

	fiberApp.Get("/", controllers.Index)

	// V1
	// UI COMPONENTS GROUP ROUTER
	components := fiberApp.Group("/v1/components")
	components.Get("/screen-choice", func(c *fiber.Ctx) error {
		return controllers.GetScreenChoices(c, embedScreens)
	})

	// INPUT VALIDATION GROUP ROUTER
	validate := fiberApp.Group("/v1/validate/json")
	validate.Post("/source", func(c *fiber.Ctx) error {
		return controllers.ValidateSrcJSONField(c, embedScreens)
	})
	validate.Post("/destination", func(c *fiber.Ctx) error {
		return controllers.ValidateDestJSONField(c, embedScreens)
	})

	fiberApp.Post("/v1/print/json", func(c *fiber.Ctx) error {
		return controllers.PrintJSON(c, embedScreens)
	})

}
