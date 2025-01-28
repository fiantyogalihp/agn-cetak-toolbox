package controllers

import (
	"embed"
	"fmt"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/utils"
	"github.com/gofiber/fiber/v2"
)

func ValidateSrcJSONField(c *fiber.Ctx, embedScreen embed.FS) error {

	screenFilename := c.FormValue("screen-choice")
	jsonInput := c.FormValue("contoh-response")

	// VALIDATE
	if screenFilename == "" {
		return utils.SendErrorResponse(c, "response-contoh", "Please select your screen!")
	}

	if jsonInput == "" {
		return utils.SendErrorResponse(c, "response-contoh", "Please input your JSON Field!")
	}

	// READ SCREEN
	screenResult, err := utils.ReadExplicitScreen(embedScreen, screenFilename+".json")
	if err != nil {
		return utils.SendErrorResponse(c, "response-contoh", err.Error())
	}

	// GET DATA
	errChan := make(chan error, 10)

	// GET UNMARSHALED RAW JSON
	rawJSON, err := utils.UnmarshalDynamicExampleJson(jsonInput)
	if err != nil {
		return utils.SendErrorResponse(c, "response-contoh", err.Error())
	}

	// CHECK EXIST FIELD
	for _, field := range screenResult.Arrange {
		if _, ok := rawJSON[field]; !ok {
			err = fmt.Errorf("invalid field, '%s' is missing", field)
			return utils.SendErrorResponse(c, "response-contoh", err.Error())
		}
	}

	utils.CheckJSONInput(errChan, rawJSON, screenResult.Required, nil)

	// ERROR HANDLE
	for err := range errChan {
		return utils.SendErrorResponse(c, "response-contoh", err.Error())
	}

	return utils.SendSuccessResponse(c, "response-contoh", "Your JSON Field is Valid!")
}

func ValidateDestJSONField(c *fiber.Ctx, embedScreen embed.FS) error {
	screenFilename := c.FormValue("screen-choice")
	jsonInput := c.FormValue("update-response")

	// VALIDATE
	if screenFilename == "" {
		return utils.SendErrorResponse(c, "response-update", "Please select your screen!")
	}
	if jsonInput == "" {
		return utils.SendErrorResponse(c, "response-update", "Please input your JSON Field!")
	}

	// READ SCREEN
	screenResult, err := utils.ReadExplicitScreen(embedScreen, screenFilename+".json")
	if err != nil {
		return utils.SendErrorResponse(c, "response-update", err.Error())
	}

	// GET DATA
	errChan := make(chan error, 10)

	// GET UNMARSHALED RAW JSON
	rawJSON, err := utils.UnmarshalDynamicExampleJson(jsonInput)
	if err != nil {
		return utils.SendErrorResponse(c, "response-update", err.Error())
	}

	checkFunc := func(data string) bool {
		return data == "pay"
	}

	// CHECK PAY FIELD
	for key := range rawJSON {
		if checkFunc(key) {
			err = fmt.Errorf("cannot process with '%s' json field", key)
			return utils.SendErrorResponse(c, "response-update", err.Error())
		}
	}

	// CHECK EXIST FIELD
	for _, field := range screenResult.Arrange {

		if checkFunc(field) {
			continue
		}

		if _, ok := rawJSON[field]; !ok {
			err = fmt.Errorf("invalid field, '%s' is missing", field)
			return utils.SendErrorResponse(c, "response-update", err.Error())
		}
	}

	utils.CheckJSONInput(errChan, rawJSON, screenResult.Required, checkFunc)

	// ERROR HANDLE
	for err := range errChan {
		return utils.SendErrorResponse(c, "response-update", err.Error())
	}

	return utils.SendSuccessResponse(c, "response-update", "Your JSON Field is Valid!")
}
