package controllers

import (
	"embed"
	"fmt"
	"log"
	"strings"

	"github.com/fiantyogalihp/agn-cetak-toolbox/utils"
	"github.com/gofiber/fiber/v2"
)

const (
	RADIO_BUTTON          string = "screen-choice"
	CONTOH_TEXTAREA       string = "contoh-response"
	ALERT_CONTOH_TEXTAREA string = "response-contoh"
	UPATE_TEXTAREA        string = "update-response"
	ALERT_UPDATE_TEXTAREA string = "response-update"
	RESULT_TEXTAREA       string = "resultTextarea"
	RESULT_BUTTON         string = "response-print"
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
	return c.Render("templates/radio_buttons", fiber.Map{
		"screens": newScreen,
	})
}
