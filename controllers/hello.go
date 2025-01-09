package controllers

import "github.com/gofiber/fiber/v2"

func HelloHandler(c *fiber.Ctx) error {
	// HTMX example: Return partial HTML
	return c.SendString("<p>Hello from Fiber!</p>")
}
