package utils

import "fmt"

func CheckJSONField(screenData *ScreenParameter, rawJson map[string]interface{}) (err error) {

	// CHECK EXIST FIELD
	for _, field := range screenData.Arrange {
		if _, ok := rawJson[field]; !ok {
			err = fmt.Errorf("invalid field, '%s' is missing", field)
			return
		}
	}

	// CHECK REQUIRED FIELD
	for _, field := range screenData.Required {

		innerField, ok := rawJson[field]

		// SKIP IF EXIST
		if ok {
			// err = fmt.Errorf("required field, '%s' is missing", field)
			// return
			continue
		}

		fmt.Println(innerField)

	}

	return
}
