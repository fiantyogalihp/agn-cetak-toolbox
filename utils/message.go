package utils

import (
	"github.com/gofiber/fiber/v2"
)

func SendErrorResponse(c *fiber.Ctx, componentID, message string) error {
	return c.Render("templates/message", fiber.Map{
		"message":      message,
		"color":        "red",
		"component_id": componentID,
	})
}

func SendSuccessResponse(c *fiber.Ctx, componentID, message string) error {
	return c.Render("templates/message", fiber.Map{
		"message":      message,
		"color":        "green",
		"component_id": componentID,
	})
}
