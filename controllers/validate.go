package controllers

import (
	"embed"
	"fmt"

	"github.com/fiantyogalihp/agn-cetak-toolbox/utils"
	"github.com/gofiber/fiber/v2"
)

func ValidateSrcJSONField(c *fiber.Ctx, embedScreen embed.FS) error {

	screenFilename := c.FormValue(RADIO_BUTTON)
	jsonInput := c.FormValue(CONTOH_TEXTAREA)

	// VALIDATE
	if screenFilename == "" {
		return utils.SendErrorResponse(c, ALERT_CONTOH_TEXTAREA, "Please select your screen!")
	}

	if jsonInput == "" {
		return utils.SendErrorResponse(c, ALERT_CONTOH_TEXTAREA, "Please input your JSON Field!")
	}

	// READ SCREEN
	screenResult, err := utils.ReadExplicitScreen(embedScreen, screenFilename+".json")
	if err != nil {
		return utils.SendErrorResponse(c, ALERT_CONTOH_TEXTAREA, err.Error())
	}

	// GET DATA
	errChan := make(chan error, 10)

	// GET UNMARSHALED RAW JSON
	rawJSON, err := utils.UnmarshalDynamicExampleJson(jsonInput)
	if err != nil {
		return utils.SendErrorResponse(c, ALERT_CONTOH_TEXTAREA, err.Error())
	}

	// CHECK EXIST FIELD
	for _, field := range screenResult.Arrange {
		if _, ok := rawJSON[field]; !ok {
			err = fmt.Errorf("invalid field, '%s' is missing", field)
			return utils.SendErrorResponse(c, ALERT_CONTOH_TEXTAREA, err.Error())
		}
	}

	utils.CheckJSONInput(errChan, rawJSON, screenResult.Required, nil)

	// ERROR HANDLE
	for err := range errChan {
		return utils.SendErrorResponse(c, ALERT_CONTOH_TEXTAREA, err.Error())
	}

	return utils.SendSuccessResponse(c, ALERT_CONTOH_TEXTAREA, "Your JSON Field is Valid!")
}

func ValidateDestJSONField(c *fiber.Ctx, embedScreen embed.FS) error {
	screenFilename := c.FormValue(RADIO_BUTTON)
	jsonInput := c.FormValue(UPATE_TEXTAREA)

	// VALIDATE
	if screenFilename == "" {
		return utils.SendErrorResponse(c, ALERT_UPDATE_TEXTAREA, "Please select your screen!")
	}
	if jsonInput == "" {
		return utils.SendErrorResponse(c, ALERT_UPDATE_TEXTAREA, "Please input your JSON Field!")
	}

	// READ SCREEN
	screenResult, err := utils.ReadExplicitScreen(embedScreen, screenFilename+".json")
	if err != nil {
		return utils.SendErrorResponse(c, ALERT_UPDATE_TEXTAREA, err.Error())
	}

	// GET DATA
	errChan := make(chan error, 10)

	// GET UNMARSHALED RAW JSON
	rawJSON, err := utils.UnmarshalDynamicExampleJson(jsonInput)
	if err != nil {
		return utils.SendErrorResponse(c, ALERT_UPDATE_TEXTAREA, err.Error())
	}

	checkFunc := func(data string) bool {
		return data == "pay"
	}

	// CHECK PAY FIELD
	for key := range rawJSON {
		if checkFunc(key) {
			err = fmt.Errorf("cannot process with '%s' json field", key)
			return utils.SendErrorResponse(c, ALERT_UPDATE_TEXTAREA, err.Error())
		}
	}

	// CHECK EXIST FIELD
	for _, field := range screenResult.Arrange {

		if checkFunc(field) {
			continue
		}

		if _, ok := rawJSON[field]; !ok {
			err = fmt.Errorf("invalid field, '%s' is missing", field)
			return utils.SendErrorResponse(c, ALERT_UPDATE_TEXTAREA, err.Error())
		}
	}

	utils.CheckJSONInput(errChan, rawJSON, screenResult.Required, checkFunc)

	// ERROR HANDLE
	for err := range errChan {
		return utils.SendErrorResponse(c, ALERT_UPDATE_TEXTAREA, err.Error())
	}

	return utils.SendSuccessResponse(c, ALERT_UPDATE_TEXTAREA, "Your JSON Field is Valid!")
}
