package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

func haveOptions(val, field, dataLocation string, errChan chan<- error) (newVal string) {
	newVal = val

	if strings.Contains(field, "-C") {
		newField := strings.ReplaceAll(field, "-C", "")

		if !strings.Contains(dataLocation, newField) {
			errChan <- fmt.Errorf("error: The field '%s' isn't same with '%s'", dataLocation, newField)
			return
		}

		newVal = newField
		return
	}

	if strings.Contains(field, "-TERBILANG") {
		// Step 1: Remove the thousand separator (.)
		cleaned := strings.ReplaceAll(dataLocation, ".", "")

		// Step 2: Replace the decimal separator (,) with a dot (.)
		cleaned = strings.ReplaceAll(cleaned, ",", ".")

		intValue, err := strconv.Atoi(cleaned)
		if err != nil {
			errChan <- err
			return
		}

		spelledOutDataLoc := ToTerbilang(int64(intValue), "RUPIAH", "upper")

		spelledOutAmount := spelledOutDataLoc + "</b>"

		newVal = spelledOutAmount
		return
	}

	if strings.Contains(field, "-N") {
		// Step 1: Remove the thousand separator (.)
		cleaned := strings.ReplaceAll(dataLocation, ".", "")

		// Step 2: Replace the decimal separator (,) with a dot (.)
		cleaned = strings.ReplaceAll(cleaned, ",", ".")

		intValue, err := strconv.Atoi(cleaned)
		if err != nil {
			errChan <- err
			return
		}

		formattedField := "Rp " + NumberFormat(float64(intValue), 0, ",", ".")

		newVal = formattedField
		return
	}

	return
}

func UnmarshalSliceToReplace(arrIndexInt int, adjustValResult map[string]interface{}, arrLocation []string, arrLocationInt []int, dataStr []byte, key string, currentLocation map[string]interface{}, errChan chan<- error) (dataLocation string) {
	switch len(arrLocation) {
	case 4:
		var nestedData [][][][]string
		err := json.Unmarshal([]byte(dataStr), &nestedData)
		if err != nil {
			log.Println(err)
			errChan <- errors.New("error: 'Data' field is not a valid JSON array")
			return
		}

		dataLocation = nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]][arrLocationInt[3]]

		if val, ok := adjustValResult[fmt.Sprint(arrIndexInt)]; ok {
			nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]][arrLocationInt[3]] = fmt.Sprint(val)

			currentLocation[key] = nestedData
		}

	case 2:
		var nestedData [][]string
		err := json.Unmarshal([]byte(dataStr), &nestedData)
		if err != nil {
			log.Println(err)
			errChan <- errors.New("error: 'Data' field is not a valid JSON array")
			return
		}

		dataLocation = nestedData[arrLocationInt[0]][arrLocationInt[1]]

		if val, ok := adjustValResult[fmt.Sprint(arrIndexInt)]; ok {
			nestedData[arrLocationInt[0]][arrLocationInt[1]] = fmt.Sprint(val)

			currentLocation[key] = nestedData
		}

	case 3:
		var nestedData [][][]string
		err := json.Unmarshal([]byte(dataStr), &nestedData)
		if err != nil {
			log.Println(err)
			errChan <- errors.New("error: 'Data' field is not a valid JSON array")
			return
		}

		dataLocation = nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]]

		if val, ok := adjustValResult[fmt.Sprint(arrIndexInt)]; ok {
			nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]] = fmt.Sprint(val)

			currentLocation[key] = nestedData
		}
	}

	return
}

func UnmarshalSlice(field string, arrLocation []string, arrLocationInt []int, dataStr []byte, errChan chan<- error) (dataLocation string) {
	switch len(arrLocation) {
	case 4:
		var nestedData [][][][]string
		err := json.Unmarshal([]byte(dataStr), &nestedData)
		if err != nil {
			errChan <- errors.New("error: 'Data' field is not a valid JSON array")
			return
		}

		dataLocation = nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]][arrLocationInt[3]]

	case 2:
		var nestedData [][]interface{}
		err := json.Unmarshal([]byte(dataStr), &nestedData)
		if err != nil {
			errChan <- errors.New("error: 'Data' field is not a valid JSON array")
			return
		}

		dataLocation = fmt.Sprint(nestedData[arrLocationInt[0]][arrLocationInt[1]])

	case 3:
		var nestedData [][][]string
		err := json.Unmarshal([]byte(dataStr), &nestedData)
		if err != nil {
			errChan <- errors.New("error: 'Data' field is not a valid JSON array")
			return
		}

		dataLocation = nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]]
	}

	// SET NEW VALUE IF HAVE OPTIONS
	dataLocation = haveOptions(dataLocation, field, dataLocation, errChan)

	return
}

func prepareProcessData(splitDataStr []string, arrLocationInt []int, arrIndexInt int, rawJSON map[string]interface{}, adjustValResult chan<- map[string]interface{}, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	field := splitDataStr[len(splitDataStr)-1]
	currentLocation := make(map[string]interface{})

	for i, v := range splitDataStr {

		arrLocation := strings.Split(v, ",")
		if len(arrLocation) < 2 {

			if len(splitDataStr) < 2 {
				if rawJSON[v] == nil {
					errChan <- fmt.Errorf("error: The field '%s' is not found", v)
					return
				}

				adjustValResult <- map[string]interface{}{fmt.Sprint(arrIndexInt): rawJSON[v]}
				return
			}

			if currentLocation[v] == nil && len(currentLocation) < 1 {
				currentLocation[v] = rawJSON[v]
				continue
			}

			previousKey := splitDataStr[i-1]
			innerData, ok := currentLocation[previousKey].(map[string]interface{})
			if !ok {
				errChan <- errors.New("error: 'Data' field JSON type is invalid")
				return
			}

			if i == len(splitDataStr)-1 {
				if innerData[v] == nil {
					errChan <- fmt.Errorf("error: The field '%s' is not found", v)
					return
				}

				adjustValResult <- map[string]interface{}{fmt.Sprint(arrIndexInt): innerData[v]}
				return
			}

			currentLocation[v] = innerData[v]
			continue
		}

		// Marshal the current location map
		previousData := splitDataStr[i-1]
		dataStr, err := json.Marshal(currentLocation[previousData])
		if err != nil {
			errChan <- err
			return
		}

		// RETURN ERROR IF SLICE NULL
		if len(arrLocationInt) < 1 {
			log.Println("data str", string(dataStr))
			errChan <- errors.New("error: slice is empty")
			return
		}

		dataLocation := UnmarshalSlice(field, arrLocation, arrLocationInt, dataStr, errChan)

		adjustValResult <- map[string]interface{}{fmt.Sprint(arrIndexInt): dataLocation}
		return
	}
}

func processReplaceData(splitDataStr []string, arrLocationInt []int, arrIndexInt int, rawJSON *map[string]interface{}, adjustValResult map[string]interface{}, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	// field := splitDataStr[len(splitDataStr)-1]
	currentLocation := make(map[string]interface{})
	resultField := make([]string, 0)

	for i, v := range splitDataStr {

		arrLocation := strings.Split(v, ",")
		if len(arrLocation) < 2 {

			if len(splitDataStr) < 2 {
				if (*rawJSON)[v] == nil {
					errChan <- fmt.Errorf("error: The field '%s' is not found", v)
					return
				}

				if val, ok := adjustValResult[fmt.Sprint(arrIndexInt)]; ok {
					currentLocation[v] = val
				}

				resultField = append(resultField, v)
				break
			}

			if currentLocation[v] == nil && len(currentLocation) < 1 {
				currentLocation[v] = (*rawJSON)[v]
				resultField = append(resultField, v)
				continue
			}

			previousKey := splitDataStr[i-1]
			currentData, ok := currentLocation[previousKey].(map[string]interface{})
			if !ok {
				errChan <- errors.New("error: 'Data' field JSON type is invalid")
				return
			}

			if i == len(splitDataStr)-1 {
				if currentData[v] == nil {
					errChan <- fmt.Errorf("error: The field '%s' is not found", v)
					return
				}

				if val, ok := adjustValResult[fmt.Sprint(arrIndexInt)]; ok {
					currentLocation[v] = val
				}

				resultField = append(resultField, v)
				break
			}

			currentLocation[v] = currentData[v]
			resultField = append(resultField, v)
			continue
		}

		// Marshal the current location map
		previousKey := splitDataStr[i-1]
		dataStr, err := json.Marshal(currentLocation[previousKey])
		if err != nil {
			errChan <- err
			return
		}

		// RETURN ERROR IF SLICE NULL
		if len(arrLocationInt) < 1 {
			errChan <- errors.New("error: slice is empty")
			return
		}

		UnmarshalSliceToReplace(arrIndexInt, adjustValResult, arrLocation, arrLocationInt, dataStr, previousKey, currentLocation, errChan)

		break
	}

	// [ "pay", "receipt" ]
	finalLevelKey := resultField[0]
	currentData := make(map[string]interface{})
	wg.Add(len(resultField))
	for i := (len(resultField) - 1); i >= 0; i-- {
		key := resultField[i]

		currentData[key] = currentLocation[key] // * Assign the value from old map to new map

		// * Skip the 1st iteration, due to no nested JSON data in those field
		if i == (len(resultField) - 1) {
			wg.Done()
			continue
		}

		// * Modify the nested JSON data want to replace
		if i < (len(resultField) - 1) {
			// * "modifiedInnerData" is the old JSON map, gonna be replaced, and make them as a new map
			modifiedInnerData, ok := currentData[key].(map[string]interface{})
			if !ok {
				errChan <- fmt.Errorf("error: hash field '%s' is invalid", key)
				return
			}
			// * assign new value from previous new JSON map to "modifiedInnerData"
			prevIdx := i + 1
			prevKey := resultField[prevIdx]
			modifiedInnerData[prevKey] = currentData[prevKey]

			currentData[key] = modifiedInnerData
		}

		wg.Done()
	}

	(*rawJSON)[finalLevelKey] = currentData[finalLevelKey]
}

func exportLocation(arrLocationChan chan<- [][]interface{}, requiredField []string) {
	defer close(arrLocationChan)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{} // Mutex to safely append to shared slice

	for i, v := range requiredField {

		splitData := strings.Split(v, ":")

		wg.Add(1)
		go func(splitData []string) {
			defer wg.Done()
			arrLocation := make([]string, 0)

			// Find the part of the string containing comma-separated integers
			for _, v := range splitData {
				splitData := strings.Split(v, ",")
				if len(splitData) > 1 {
					arrLocation = splitData
					break
				}
			}

			splitDataInterface := make([]interface{}, 0)
			for _, v := range splitData {
				splitDataInterface = append(splitDataInterface, v)
			}

			// Convert strings to integers
			arrLocationInt := make([]interface{}, 0)
			for _, v := range arrLocation {
				dataInt, err := strconv.Atoi(v)
				if err != nil {
					log.Println(err)
					return
				}
				arrLocationInt = append(arrLocationInt, dataInt)
			}

			arrIndexInt := []interface{}{i}

			// Send result to the channel (safely)
			mu.Lock()
			defer mu.Unlock()
			arrLocationChan <- [][]interface{}{splitDataInterface, arrLocationInt, arrIndexInt}
		}(splitData)
	}

	wg.Wait()
}

func PrepareJSONInput(errChan chan<- error, rawJSON map[string]interface{}, fieldValue []string, adjustValResult chan<- map[string]interface{}) {

	arrLocationValChan := make(chan [][]interface{})
	go exportLocation(arrLocationValChan, fieldValue)

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	for data := range arrLocationValChan {

		splitDataInterface := data[0]
		arrLocationInterface := data[1]
		arrIndexInterface := data[2]

		splitDataStrChan := make(chan []string, 10)
		arrLocationIntChan := make(chan []int, 10)
		arrIndexIntChan := make(chan []int, 10)

		go convertSliceStr(splitDataInterface, splitDataStrChan, errChan, &wg, &mu)
		go convertSliceInt(arrLocationInterface, arrLocationIntChan, errChan, &wg, &mu)
		go convertSliceInt(arrIndexInterface, arrIndexIntChan, errChan, &wg, &mu)

		wg.Add(1)
		go func(rawJSON map[string]interface{}, adjustValResult chan<- map[string]interface{}, splitDataStrChan <-chan []string, arrLocationIntChan, arrIndexIntChan <-chan []int) {
			defer wg.Done()

			splitDataStr := <-splitDataStrChan
			arrLocationInt := <-arrLocationIntChan
			arrIndexInt := <-arrIndexIntChan

			mu.Lock()
			defer mu.Unlock()
			prepareProcessData(splitDataStr, arrLocationInt, arrIndexInt[0], rawJSON, adjustValResult, errChan, &wg)
		}(rawJSON, adjustValResult, splitDataStrChan, arrLocationIntChan, arrIndexIntChan)

	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(adjustValResult)
		close(errChan)
	}()
}

func PrintJSONInput(errChan chan<- error, rawJSON map[string]interface{}, fieldKey []string, adjustValResult map[string]interface{}, resultJSON *map[string]interface{}) {
	defer close(errChan)

	arrLocationKeyChan := make(chan [][]interface{})
	go exportLocation(arrLocationKeyChan, fieldKey)

	wg := sync.WaitGroup{}
	mu := &sync.Mutex{}

	resultMapChan := make(chan map[string]interface{})
	newRawJSON := rawJSON

	// Read from the channel
	for data := range arrLocationKeyChan {

		splitDataInterface := data[0]
		arrLocationInterface := data[1]
		arrIndexInterface := data[2]

		splitDataStrChan := make(chan []string, 10)
		arrLocationIntChan := make(chan []int, 10)
		arrIndexIntChan := make(chan []int, 10)

		go convertSliceStr(splitDataInterface, splitDataStrChan, errChan, &wg, mu)
		go convertSliceInt(arrLocationInterface, arrLocationIntChan, errChan, &wg, mu)
		go convertSliceInt(arrIndexInterface, arrIndexIntChan, errChan, &wg, mu)

		wg.Add(1)
		go func(adjustValResult map[string]interface{}, splitDataStrChan <-chan []string, arrLocationIntChan, arrIndexIntChan <-chan []int, resultMapChan chan<- map[string]interface{}) {
			defer wg.Done()

			splitDataStr := <-splitDataStrChan
			arrLocationInt := <-arrLocationIntChan
			arrIndexInt := <-arrIndexIntChan

			mu.Lock()
			defer mu.Unlock()

			processReplaceData(splitDataStr, arrLocationInt, arrIndexInt[0], &newRawJSON, adjustValResult, errChan, &wg)

		}(adjustValResult, splitDataStrChan, arrLocationIntChan, arrIndexIntChan, resultMapChan)

	}

	// Wait for all goroutines to finish
	// go func() {
	wg.Wait()
	// }()

	log.Println("new raw json", newRawJSON)

	*resultJSON = newRawJSON
}

func MarshalFinalResult(resultJSON map[string]interface{}) (result string, err error) {

	jsonInqData, err := json.Marshal(resultJSON["inq"])
	if err != nil {
		return
	}

	jsonPayData, err := json.Marshal(resultJSON["pay"])
	if err != nil {
		return
	}

	resultJSON["inq"] = string(jsonInqData)
	resultJSON["pay"] = string(jsonPayData)

	jsonResultData, err := json.Marshal(resultJSON)
	if err != nil {
		return
	}

	result = string(jsonResultData)
	return
}
