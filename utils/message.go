package utils

import "github.com/gofiber/fiber/v2"

func SendErrorResponse(c *fiber.Ctx, message string) error {
	return c.Render("templates/message", fiber.Map{
		"message": message,
		"color":   "red",
	})
}

func SendSuccessResponse(c *fiber.Ctx, message string) error {
	return c.Render("templates/message", fiber.Map{
		"message": message,
		"color":   "green",
	})
}
