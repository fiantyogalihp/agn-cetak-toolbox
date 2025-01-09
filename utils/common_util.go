package utils

import (
	"encoding/json"
	"errors"
	"log"
)

func UnmarshalDynamicExampleJson(rawJson string) (result map[string]interface{}, err error) {
	// Unmarshal the outer JSON into a map
	outer := make(map[string]interface{}, 0)
	err = json.Unmarshal([]byte(rawJson), &outer)
	if err != nil {
		log.Fatal("Error unmarshaling outer JSON:", err)
		return
	}

	// Unmarshal the inner JSON into a map
	for key, value := range outer {
		// Check if "Data" is a valid JSON object
		dataStr, err := checkTypeDataOfJson(value)
		if err != nil {
			if err.Error() == "CONTINUE" {
				continue
			}

			log.Fatal(err)
			return result, err
		}

		inner := make(map[string]interface{}, 0)
		err = json.Unmarshal([]byte(dataStr), &inner)
		if err != nil {
			log.Fatal("Error unmarshaling inner JSON:", err)
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

	case bool, int, float64:
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
