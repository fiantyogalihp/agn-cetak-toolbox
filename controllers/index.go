package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	// Return the HTML response
	return c.Render("templates/index", fiber.Map{
		"Title": "HTMX + Fiber Quickstart",
	})
}
func Screen(c *fiber.Ctx) error {
	// Return the HTML response
	return c.Render("templates/screen", nil)
}
