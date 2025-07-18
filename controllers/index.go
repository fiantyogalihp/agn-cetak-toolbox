package controllers

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	title := "AGN Cetak Toolbox"

	visited := c.Query("visited")
	if visited != "" && visited == "true" {
		return c.Render("templates/screen", fiber.Map{
			"Title": title,
		})
	}

	return c.Render("templates/index", fiber.Map{
		"Title": title,
	})
}
