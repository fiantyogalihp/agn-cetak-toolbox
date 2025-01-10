package controllers

import (
	"embed"
	"log"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/utils"
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx, embedScreen embed.FS) error {

	screenResult, err := utils.ReadScreen(embedScreen)
	if err != nil {
		log.Println(err)
		return err
	}

	newScreen := make([]map[string]interface{}, 0)
	for _, screen := range screenResult {
		if _, ok := screen["filename"]; ok {
			newScreen = append(newScreen, screen)
		}
	}

	// Generate HTML dynamically
	// var html string
	// for _, item := range items {
	// html += fmt.Sprintf("<input type=\"radio\" name=\"%s\">", item.Name)
	// }

	// Return the HTML response
	return c.Render("templates/index", fiber.Map{
		"Title":  "HTMX + Fiber Quickstart",
		"Screen": newScreen,
	})
}
