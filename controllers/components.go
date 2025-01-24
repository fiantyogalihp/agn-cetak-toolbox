package controllers

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/utils"
	"github.com/gofiber/fiber/v2"
)

func GetScreenChoices(c *fiber.Ctx, embedScreen embed.FS) error {
	screenResult, err := utils.ReadScreen(embedScreen)
	if err != nil {
		log.Println(err)
		return err
	}

	newScreen := make([]map[string]string, 0)
	for _, screen := range screenResult {
		if _, ok := screen["filename"]; ok {
			value := strings.ReplaceAll(fmt.Sprint(screen["filename"]), ".json", "")

			newScreen = append(newScreen, map[string]string{
				"value": value,
				"label": fmt.Sprint(screen["screen_name"]),
			})
		}
	}

	// Return the HTML response
	return c.Render("templates/radio-buttons", fiber.Map{
		"screens": newScreen,
	})
}
