package controllers

import (
	"embed"
	"fmt"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/utils"
	"github.com/gofiber/fiber/v2"
)

func ValidateJsonField(c *fiber.Ctx, embedScreen embed.FS) error {

	screenFilename := c.FormValue("screen-choice") + ".json"

	// READ SCREEN
	screeenResult, err := utils.ReadExplicitScreen(embedScreen, screenFilename)
	if err != nil {
		errMessage := fmt.Sprintf("<div style='color: red;'>%s</div>", err.Error())
		return c.SendString(errMessage)
	}

	// GET DATA
	jsonInput := c.FormValue("contoh-response")
	errChan := make(chan error, 10)

	// VALIDATE
	utils.CheckExampleJSONInput(errChan, jsonInput, screeenResult.Required)

	// ERROR HANDLE
	for err := range errChan {
		errMessage := fmt.Sprintf("<div style='color: red;'>%s</div>", err.Error())
		return c.SendString(errMessage)
	}

	return c.SendString("<div style='color: green;'>Your JSON Field is Valid!</div>")
}
