package controllers

import (
	"embed"
	"fmt"
	"log"
	"sync"

	"github.com/fiantyogalihp/agn-cetak-toolbox/utils"
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
	errPrepareChan := make(chan error, 10)
	errPrintChan := make(chan error, 10)
	errChan := make(chan error, 10)
	mu := sync.Mutex{}

	// PREPARE INPUT DATA
	rawExampleJSONChan := make(chan map[string]interface{})
	rawUpdateJSONChan := make(chan map[string]interface{})
	go func() {
		defer close(rawExampleJSONChan)
		defer close(rawUpdateJSONChan)
		defer close(errChan)

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

	for err := range errChan {
		return utils.SendErrorResponse(c, "response-print", err.Error())
	}

	rawUpdateJSON["pay"] = rawExampleJSON["pay"]

	fieldKey := make([]string, 0)
	fieldValue := make([]string, 0)
	for k, v := range screenResult.Adjustment {
		fieldKey = append(fieldKey, k)
		fieldValue = append(fieldValue, v)
	}

	log.Println(fieldKey, fieldValue)

	// PREPARE SOURCE DATA
	adjustValResult := make(chan map[string]interface{})
	utils.PrepareJSONInput(errPrepareChan, rawUpdateJSON, fieldValue, adjustValResult)

	adjustValMap := make(map[string]interface{})
	for data := range adjustValResult {

		counterKey := 0
		key := utils.FindKey(counterKey, len(fieldKey), data)
		if key > len(fieldKey) {
			errPrepareChan <- fmt.Errorf("key '%d' not found", key)
			continue
		}

		keyStr := fmt.Sprint(key)

		adjustValMap[keyStr] = data[keyStr]
	}

	for err := range errPrepareChan {
		return utils.SendErrorResponse(c, "response-print", err.Error())
	}

	// REPLACE RESULT DATA
	resultJSON := make(map[string]interface{})
	utils.PrintJSONInput(errPrintChan, rawUpdateJSON, fieldKey, adjustValMap, &resultJSON)

	for err := range errPrintChan {
		return utils.SendErrorResponse(c, "response-print", err.Error())
	}

	// MARSHAL TO FIRST FORMAT
	jsonResultData, err := utils.MarshalFinalResult(resultJSON)
	if err != nil {
		return utils.SendErrorResponse(c, "response-print", err.Error())
	}

	// ADD OPTION DATA
	arrangedArr := screenResult.Arrange

	return c.Render("templates/textarea_result", fiber.Map{
		"result":  jsonResultData,
		"arrange": arrangedArr,
	})
}
