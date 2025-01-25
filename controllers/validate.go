package controllers

import (
	"embed"
	"fmt"
	"log"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/utils"
	"github.com/gofiber/fiber/v2"
)

func ValidateJSONField(c *fiber.Ctx, embedScreen embed.FS) error {

	screenFilename := c.FormValue("screen-choice")

	// VALIDATE
	if screenFilename == "" {
		return utils.SendErrorResponse(c, "Please select your screen!")
	}

	screenFilename += ".json"

	// READ SCREEN
	screeenResult, err := utils.ReadExplicitScreen(embedScreen, screenFilename)
	if err != nil {
		return utils.SendErrorResponse(c, err.Error())
	}

	// GET DATA
	jsonInput := c.FormValue("contoh-response")
	errChan := make(chan error, 10)

	// VALIDATE
	if jsonInput == "" {
		return utils.SendErrorResponse(c, "Please input your JSON Field!")
	}

	rawJSONChan := make(chan map[string]interface{})
	go func() {
		defer close(rawJSONChan)

		result, err := utils.UnmarshalDynamicExampleJson(jsonInput)
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		rawJSONChan <- result

	}()
	rawJSON := <-rawJSONChan

	// CHECK EXIST FIELD
	for _, field := range screeenResult.Arrange {
		if _, ok := rawJSON[field]; !ok {
			err = fmt.Errorf("invalid field, '%s' is missing", field)
			return utils.SendErrorResponse(c, err.Error())
		}
	}

	utils.CheckExampleJSONInput(errChan, jsonInput, screeenResult.Required)

	// ERROR HANDLE
	for err := range errChan {
		return utils.SendErrorResponse(c, err.Error())
	}

	return utils.SendSuccessResponse(c, "Your JSON Field is Valid!")
}
