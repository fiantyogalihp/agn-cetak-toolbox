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

func processSplitData(splitDataStr []string, arrLocationInt []int, rawJSON map[string]interface{}, currentLocation map[string]interface{}, errChan chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	field := splitDataStr[len(splitDataStr)-1]

	for i, v := range splitDataStr {
		// fmt.Println(i, "ini value", v, len(splitDataStr))

		arrLocation := strings.Split(v, ",")
		if len(arrLocation) < 2 {

			if len(splitDataStr) < 2 {
				if v != field {
					errChan <- fmt.Errorf("error: The field '%s' isn't same with '%s'", v, field)
					return
				}

				fmt.Println("Final result:", v, field)
				return
			}

			// fmt.Println(i, "ini value nya:", v)
			// fmt.Println(i, "ini json nya:", rawJSON[v])
			if currentLocation[v] == nil && len(currentLocation) < 1 {
				currentLocation[v] = rawJSON[v]
				continue
			}

			previousData := splitDataStr[i-1]

			// fmt.Println(i, "previous location:", currentLocation[previousData])
			// fmt.Println(i, reflect.TypeOf(currentLocation[previousData]))
			innerData, ok := currentLocation[previousData].(map[string]interface{})
			if !ok {
				errChan <- errors.New("error: 'Data' field JSON type is invalid")
				return
			}

			if i == len(splitDataStr)-1 {
				if v != field {
					errChan <- fmt.Errorf("error: The field '%s' isn't same with '%s'", innerData[v], field)
					return
				}
				fmt.Println("Final result:", v, field)
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

		switch len(arrLocation) {
		case 4:
			var nestedData [][][][]string
			err = json.Unmarshal([]byte(dataStr), &nestedData)
			if err != nil {
				errChan <- errors.New("error: 'Data' field is not a valid JSON array")
				return
			}

			// fmt.Println("ini nested data", nestedData, arrLocationInt)
			// fmt.Println("ini array", arrLocationInt[0], arrLocationInt[1], arrLocationInt[2], arrLocationInt[3])
			dataLocation := nestedData[arrLocationInt[0]][arrLocationInt[1]][arrLocationInt[2]][arrLocationInt[3]]
			if strings.Contains(field, "-C") {
				fieldNew := strings.ReplaceAll(field, "-C", "")

				if !strings.Contains(dataLocation, fieldNew) {
					errChan <- fmt.Errorf("error: The field '%s' isn't same with '%s'", dataLocation, fieldNew)
					return
				}

				fmt.Println("Final result: ", dataLocation, field)
				return
			}

			if dataLocation != field {
				errChan <- fmt.Errorf("error: The field '%s' isn't same with '%s'", dataLocation, field)
				return
			}

			fmt.Println("Final result: ", dataLocation, field)
			return
		}
	}
}

func convertSliceStr(slice []interface{}, resultStrChan chan<- []string, errChan chan<- error, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	wg.Add(1)

	defer close(resultStrChan)

	dataJson, err := json.Marshal(slice)
	if err != nil {
		errChan <- err
		return
	}

	// Convert the interface slice to a slice of strings
	splitDataStr := make([]string, len(slice))
	err = json.Unmarshal(dataJson, &splitDataStr)
	if err != nil {
		errChan <- err
		return
	}

	mu.Lock()
	resultStrChan <- splitDataStr
	mu.Unlock()

}

func convertSliceInt(slice []interface{}, resultIntChan chan<- []int, errChan chan<- error, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()
	wg.Add(1)

	defer close(resultIntChan)

	dataJson, err := json.Marshal(slice)
	if err != nil {
		errChan <- err
		return
	}

	// Convert the interface slice to a slice of int
	splitDataInt := make([]int, len(slice))
	err = json.Unmarshal(dataJson, &splitDataInt)
	if err != nil {
		errChan <- err
		return
	}

	mu.Lock()
	resultIntChan <- splitDataInt
	mu.Unlock()

}

func CheckJSONInput(errChan chan<- error, rawJSON map[string]interface{}, requiredField []string, shouldSkip func(data string) bool) {

	arrLocationChan := make(chan [][]interface{})
	go func() {
		defer close(arrLocationChan)

		wg := sync.WaitGroup{}
		mu := sync.Mutex{} // Mutex to safely append to shared slice

		for _, v := range requiredField {

			splitData := strings.Split(v, ":")

			// If the check function is provided and returns true, skip this iteration
			// CONTINUE 'PAY' FIELD CHECK
			if shouldSkip != nil && shouldSkip(splitData[0]) {
				continue
			}

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

				// Send result to the channel (safely)
				mu.Lock()
				arrLocationChan <- [][]interface{}{splitDataInterface, arrLocationInt}
				mu.Unlock()
			}(splitData)
		}

		wg.Wait()
	}()

	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	// fmt.Println(rawJSON)

	// Read from the channel
	for data := range arrLocationChan {

		splitDataInterface := data[0]
		arrLocationInterface := data[1]

		splitDataStrChan := make(chan []string, 10)
		arrLocationIntChan := make(chan []int, 10)

		go convertSliceStr(splitDataInterface, splitDataStrChan, errChan, &wg, &mu)
		go convertSliceInt(arrLocationInterface, arrLocationIntChan, errChan, &wg, &mu)

		wg.Add(1)
		go func(rawJSON map[string]interface{}) {
			defer wg.Done()

			splitDataStr := <-splitDataStrChan
			arrLocationInt := <-arrLocationIntChan

			// field and currentLocation
			// field := splitDataStr[len(splitDataStr)-1]
			currentLocation := make(map[string]interface{})

			// fmt.Println(arrLocationInt)

			mu.Lock()
			processSplitData(splitDataStr, arrLocationInt, rawJSON, currentLocation, errChan, &wg)
			mu.Unlock()
		}(rawJSON)

	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(errChan)
	}()

}
