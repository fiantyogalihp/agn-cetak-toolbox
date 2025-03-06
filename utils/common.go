package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func UnmarshalDynamicExampleJson(rawJson string) (result map[string]interface{}, err error) {
	// Unmarshal the outer JSON into a map
	decoder := json.NewDecoder(strings.NewReader(rawJson))
	decoder.UseNumber() // Prevents conversion to float64

	outer := make(map[string]interface{})
	if err := decoder.Decode(&outer); err != nil {
		log.Println(err)
		return result, err
	}

	// Unmarshal the inner JSON into a map
	for key, value := range outer {
		// Check if "Data" is a valid JSON object
		dataStr, err := checkTypeDataOfJson(value)
		if err != nil {
			if err.Error() == "CONTINUE" {
				continue
			}

			log.Println(err)
			return result, err
		}

		decoder := json.NewDecoder(strings.NewReader(dataStr))
		decoder.UseNumber() // Prevents conversion to float64

		inner := make(map[string]interface{})
		if err := decoder.Decode(&inner); err != nil {
			log.Println(err)
			return result, err
		}

		outer[key] = inner
	}

	result = outer

	return
}

func checkTypeDataOfJson(dataValue interface{}) (result string, err error) {
	// Check if "Data" is a valid JSON object
	switch dataValue.(type) {
	case string:
		// "Data" is a string - try to unmarshal it into a JSON object or array
		dataStr, _ := dataValue.(string)

		var nestedData interface{}
		err = json.Unmarshal([]byte(dataStr), &nestedData)
		if err != nil {
			// ! error: 'Data' field is not a valid JSON object or array
			err = errors.New("CONTINUE")
			return
		}

		// Check if the unmarshaled data is an object or an array
		switch nestedData.(type) {
		case map[string]interface{}, []interface{}: // JSON object or array
			result = dataStr
			return
		default:
			// ! error: 'Data' field is not a valid JSON object or array
			err = errors.New("CONTINUE")
			return
		}

	case bool, int, float64, json.Number:
		// ! error: 'Data' field is not a nested JSON object
		err = errors.New("CONTINUE")
		return
	default:
		// Handle case where "Data" is expected to be a JSON string
		result, ok := dataValue.(string)
		if !ok {
			err = errors.New("error: 'Data' field type is invalid")
			return result, err
		}
	}

	return
}

func NumberFormat(number float64, decimals uint, decPoint, thousandsSep string) string {
	neg := false
	if number < 0 {
		number = -number
		neg = true
	}
	dec := int(decimals)
	// Will round off
	str := fmt.Sprintf("%."+strconv.Itoa(dec)+"F", number)
	prefix, suffix := "", ""
	if dec > 0 {
		prefix = str[:len(str)-(dec+1)]
		suffix = str[len(str)-dec:]
	} else {
		prefix = str
	}
	sep := []byte(thousandsSep)
	n, l1, l2 := 0, len(prefix), len(sep)
	// thousands sep num
	csep := (l1 - 1) / 3
	tmp := make([]byte, l2*csep+l1)
	pos := len(tmp) - 1
	for i := l1 - 1; i >= 0; i, n, pos = i-1, n+1, pos-1 {
		if l2 > 0 && n > 0 && n%3 == 0 {
			for j := range sep {
				tmp[pos] = sep[l2-j-1]
				pos--
			}
		}
		tmp[pos] = prefix[i]
	}
	s := string(tmp)
	if dec > 0 {
		s += decPoint + suffix
	}
	if neg {
		s = "-" + s
	}

	return s
}

// Description:
// This function recursively searches for a key in the `target` map starting from `counterKey` and increments the key until it finds a key that does not exist in the map or exceeds the `max` value.
//
// Parameters:
// - `counterKey` (int): The starting key value for the search similar as index.
// - `max` (int): The maximum value for the key search.
// - `target` (map[string]string): The map in which to search for the key.
//
// Returns:
// - `string`: The key that is either not found in the map or exceeds the `max` value.
//
// Example:
// ```go
//
//	data := map[string]string{
//	  "12": "fiant",
//	}
//
//	counterKey := 0
//	key := findKey(counterKey, 10, data)
//	fmt.Println(key) // Output will be "7"
//
// ```
func FindKey(counterKey int, max int, target map[string]interface{}) int {
	key := counterKey

	if counterKey > max {
		return key
	}

	if _, ok := target[fmt.Sprint(counterKey)]; !ok {
		key = FindKey(counterKey+1, max, target)
	}

	return key
}
