package utils

import "github.com/gofiber/fiber/v2"

func SendErrorResponse(c *fiber.Ctx, templates, message string) error {
	return c.Render(templates, fiber.Map{
		"error": message,
	})
}
