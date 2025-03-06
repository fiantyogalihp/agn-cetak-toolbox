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

		// fmt.Println("nested data", nestedData)

		// fmt.Println("ini nested data", nestedData, arrLocationInt)
		// fmt.Println("ini array", arrLocationInt[0], arrLocationInt[1], arrLocationInt[2], arrLocationInt[3])
		dataLocation = nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]][arrLocationInt[3]]

		if val, ok := adjustValResult[fmt.Sprint(arrIndexInt)]; ok {
			// newNestedData := make([][][][]string, len(nestedData))
			// copy(newNestedData, nestedData[:])

			// newVal := haveOptions(val, field, dataLocation, errChan)

			// fmt.Println("new nested data", newNestedData)

			nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]][arrLocationInt[3]] = fmt.Sprint(val)

			currentLocation[key] = nestedData
			// *resultField = append(*resultField, key)
			// inputList.PushBack(map[string]interface{}{key: newNestedData})
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
			// newNestedData := make([][]string, len(nestedData))
			// copy(newNestedData, nestedData[:])

			// newVal := haveOptions(val, field, dataLocation, errChan)

			nestedData[arrLocationInt[0]][arrLocationInt[1]] = fmt.Sprint(val)

			currentLocation[key] = nestedData
			// *resultField = append(*resultField, key)
			// inputList.PushBack(map[string]interface{}{key: newNestedData})
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
			// newNestedData := make([][][]string, len(nestedData))
			// copy(newNestedData, nestedData[:])

			// newVal := haveOptions(val, field, dataLocation, errChan)

			nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]] = fmt.Sprint(val)

			currentLocation[key] = nestedData
			// *resultField = append(*resultField, key)
			// inputList.PushBack(map[string]interface{}{key: newNestedData})
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

		// fmt.Println("ini nested data", nestedData, arrLocationInt)
		// fmt.Println("ini array", arrLocationInt[0], arrLocationInt[1], arrLocationInt[2], arrLocationInt[3])
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

	// fmt.Println("ini value", splitDataStr, len(splitDataStr))
	for i, v := range splitDataStr {
		// fmt.Println(i, "ini value", v, len(splitDataStr))

		arrLocation := strings.Split(v, ",")
		if len(arrLocation) < 2 {

			if len(splitDataStr) < 2 {
				if rawJSON[v] == nil {
					errChan <- fmt.Errorf("error: The field '%s' is not found", v)
					return
				}

				// fmt.Println("masuk sini", arrIndexInt, rawJSON[v])
				adjustValResult <- map[string]interface{}{fmt.Sprint(arrIndexInt): rawJSON[v]}
				return
			}

			// fmt.Println(i, "ini json nya:", rawJSON[v])
			// fmt.Println(i, "ini value nya:", v)
			if currentLocation[v] == nil && len(currentLocation) < 1 {
				currentLocation[v] = rawJSON[v]
				continue
			}

			previousKey := splitDataStr[i-1]
			// fmt.Println(i, "previous location:", currentLocation[previousKey])
			// fmt.Println(i, reflect.TypeOf(currentLocation[previousKey]))
			innerData, ok := currentLocation[previousKey].(map[string]interface{})
			if !ok {
				errChan <- errors.New("error: 'Data' field JSON type is invalid")
				return
			}

			// fmt.Println("ini inner data", innerData[v])

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

		// fmt.Println("current data", currentLocation)
		// fmt.Println("previous data", currentLocation[previousData])

		dataStr, err := json.Marshal(currentLocation[previousData])
		if err != nil {
			errChan <- err
			return
		}

		// RETURN ERROR IF SLICE NULL
		if len(arrLocationInt) < 1 {
			fmt.Println("datastr", string(dataStr))
			errChan <- errors.New("error: slice is empty")
			return
		}

		dataLocation := UnmarshalSlice(field, arrLocation, arrLocationInt, dataStr, errChan)

		// if dataLocation != field {
		// 	errChan <- fmt.Errorf("error: The field '%s' isn't same with '%s'", dataLocation, field)
		// 	return
		// }
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
		// fmt.Println(i, "ini value", v, len(splitDataStr))

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
				// fmt.Println("________________________________\n", v, resultField)
				// fmt.Println("masuk sini", arrIndexInt, rawJSON[v])
				// adjustKeyResult <- map[string]string{fmt.Sprint(arrIndexInt): fmt.Sprint(rawJSON[v])}
				break
			}

			// fmt.Println(i, "ini json nya:", rawJSON[v])
			// fmt.Println(i, "ini value nya:", v)
			if currentLocation[v] == nil && len(currentLocation) < 1 {
				currentLocation[v] = (*rawJSON)[v]
				resultField = append(resultField, v)
				// fmt.Println("________________________________\n", v, resultField)
				// resultList.PushBack(rawJSON[v])
				continue
			}

			previousKey := splitDataStr[i-1]
			// fmt.Println(i, "previous location:", currentLocation[previousKey])
			// fmt.Println(i, reflect.TypeOf(currentLocation[previousKey]))
			currentData, ok := currentLocation[previousKey].(map[string]interface{})
			if !ok {
				errChan <- errors.New("error: 'Data' field JSON type is invalid")
				return
			}

			// fmt.Println("ini inner data", currentData[v])

			if i == len(splitDataStr)-1 {
				if currentData[v] == nil {
					errChan <- fmt.Errorf("error: The field '%s' is not found", v)
					return
				}

				if val, ok := adjustValResult[fmt.Sprint(arrIndexInt)]; ok {
					currentLocation[v] = val
				}

				resultField = append(resultField, v)
				// fmt.Println("________________________________\n", v, resultField)
				// adjustKeyResult <- map[string]string{fmt.Sprint(arrIndexInt): fmt.Sprint(currentData[v])}
				break
			}

			currentLocation[v] = currentData[v]
			resultField = append(resultField, v)
			// fmt.Println("________________________________\n", v, resultField)
			// resultList.PushBack(currentData[v])
			continue
		}

		// Marshal the current location map
		previousKey := splitDataStr[i-1]
		// fmt.Println("current location", currentLocation)
		// fmt.Println("current location data str", currentLocation[previousKey], previousKey)
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

		// fmt.Println("________________________________\n", previousKey, dataLocation)
		// currentLocation[previousKey] = dataLocation
		// adjustKeyResult <- map[string]string{fmt.Sprint(arrIndexInt): dataLocation}
		break
	}

	fmt.Println(resultField)

	// [ "pay", "receipt" ]
	finalLevelKey := resultField[0]
	currentData := make(map[string]interface{})
	wg.Add(len(resultField))
	for i := (len(resultField) - 1); i >= 0; i-- {
		key := resultField[i]

		// fmt.Println("current data 0", key, currentData)

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

			// fmt.Println(modifiedInnerData)
			// fmt.Println(key, prevIdx, resultField)

			// * assign new value from previous new JSON map to "modifiedInnerData"
			prevIdx := i + 1
			prevKey := resultField[prevIdx]
			modifiedInnerData[prevKey] = currentData[prevKey]

			currentData[key] = modifiedInnerData
		}

		wg.Done()

		// if nextIdx := i - 1; nextIdx >= 0 {
		// 	currentData[key] = currentLocation[key]

		// 	fmt.Println("current data 1", currentData[key], currentLocation[key])

		// 	innerData, ok := currentData[key].(map[string]interface{})
		// 	if !ok {
		// 		errChan <- errors.New("error: 'Data' field JSON type is invalid")
		// 		return
		// 	}

		// 	innerData[key] = currentLocation[key]

		// 	// currentData[key] = currentLocation[key] // * (*new) receipt = newValue
		// 	continue
		// }

		// prevIdx := i + 1
		// fmt.Println(key, prevIdx, resultField)
		// prevKey := resultField[prevIdx]

		// fmt.Println("________________________________", currentData[prevKey])

		// currentData[key] = currentData[prevKey] // * (*new) pay = (*new) receipt
	}

	// fmt.Println("current data", map[string]interface{}{finalLevelKey: currentData[finalLevelKey]})
	// fmt.Println("sended", finalLevelKey)
	// resultMapChan <- map[string]interface{}{finalLevelKey: currentData[finalLevelKey]}
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
		// fmt.Println("data", data)

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

			// fmt.Println(splitDataStr, arrLocationInt, arrIndexInt)
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
	// cond := sync.NewCond(mu)

	resultMapChan := make(chan map[string]interface{})
	newRawJSON := rawJSON
	// var finalLevelKey string

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

		// time.Sleep(10 * time.Millisecond)

		wg.Add(1)
		go func(adjustValResult map[string]interface{}, splitDataStrChan <-chan []string, arrLocationIntChan, arrIndexIntChan <-chan []int, resultMapChan chan<- map[string]interface{}) {
			defer wg.Done()

			splitDataStr := <-splitDataStrChan
			arrLocationInt := <-arrLocationIntChan
			arrIndexInt := <-arrIndexIntChan

			mu.Lock()
			defer mu.Unlock()
			// cond.L.Lock()
			// defer cond.L.Unlock()
			// finalLevelKey = splitDataStr[0]

			processReplaceData(splitDataStr, arrLocationInt, arrIndexInt[0], &newRawJSON, adjustValResult, errChan, &wg)

			// cond.Wait()
			// newRawJSON[finalLevelKey] = <-resultMapChan
		}(adjustValResult, splitDataStrChan, arrLocationIntChan, arrIndexIntChan, resultMapChan)

	}

	// wg.Add(1)
	// go func() {
	// 	defer close(resultMapChan)
	// 	defer wg.Done()

	// 	}()
	// for newData := range resultMapChan {

	// 	wg.Add(1)
	// 	go func() {
	// 		defer wg.Done()
	// 		cond.L.Lock()
	// 		newRawJSON[finalLevelKey] = newData[finalLevelKey]
	// 		fmt.Println("success")
	// 		cond.L.Unlock()
	// 		cond.Signal()
	// 	}()

	// 	// fmt.Println("________________________________\n", newRawJSON[finalLevelKey])
	// 	// fmt.Println("________________________________\n", newData[finalLevelKey])
	// 	// fmt.Println("________________________________\n", resultJSON)
	// }

	// close(resultMapChan)
	// Wait for all goroutines to finish
	// go func() {
	wg.Wait()
	// }()

	fmt.Println("new raw json", newRawJSON)

	*resultJSON = newRawJSON
	// resultJSON <- newRawJSON
	// close(resultJSON)
}
