package controllers

import (
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/fiantyogalihp/dynamic-json-parsing-struct/utils"
	"github.com/gofiber/fiber/v2"
)

func PrintJSON(c *fiber.Ctx, embedScreen embed.FS) error {

	screenFilename := c.FormValue("screen-choice")
	exampleJSONInput := c.FormValue("contoh-response")
	updateJSONInput := c.FormValue("update-response")

	// VALIDATE
	if screenFilename == "" {
		return utils.SendErrorResponse(c, "print-response", "Please select your screen!")
	}

	if exampleJSONInput == "" || updateJSONInput == "" {
		return utils.SendErrorResponse(c, "print-response", "Please input your JSON Field!")
	}

	// READ SCREEN
	screenResult, err := utils.ReadExplicitScreen(embedScreen, screenFilename+".json")
	if err != nil {
		return utils.SendErrorResponse(c, "print-response", err.Error())
	}

	// GET DATA
	errChan := make(chan error, 10)
	mu := sync.Mutex{}

	rawExampleJSONChan := make(chan map[string]interface{})
	rawUpdateJSONChan := make(chan map[string]interface{})
	go func() {
		defer close(rawExampleJSONChan)
		defer close(rawUpdateJSONChan)

		resultExample, err := utils.UnmarshalDynamicExampleJson(exampleJSONInput)
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		resultUpdate, err := utils.UnmarshalDynamicExampleJson(updateJSONInput)
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		mu.Lock()
		rawExampleJSONChan <- resultExample
		rawUpdateJSONChan <- resultUpdate
		mu.Unlock()

	}()
	rawExampleJSON := <-rawExampleJSONChan
	rawUpdateJSON := <-rawUpdateJSONChan

	rawUpdateJSON["pay"] = rawExampleJSON["pay"]

	jsonData, err := json.Marshal(rawUpdateJSON)
	if err != nil {
		return utils.SendErrorResponse(c, "response-print", err.Error())
	}

	fmt.Println(screenResult)

	return c.Render("templates/textarea_result", fiber.Map{
		"result": string(jsonData),
	})
}
