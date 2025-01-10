package utils

import (
	"embed"
	"encoding/json"
	"fmt"
	"path/filepath"
)

func ReadScreen(screenData embed.FS) (result []map[string]interface{}, err error) {

	// Get a list of all files in the "screens" directory
	files, err := screenData.ReadDir("screens")
	if err != nil {
		err = fmt.Errorf("failed to read screens directory: %v", err)
		return
	}

	// Loop through the files and print their names and contents
	result = make([]map[string]interface{}, 0)
	for _, file := range files {
		// Skip directories (if any)
		if file.IsDir() {
			continue
		}

		// Read the file's content
		content, err := screenData.ReadFile(filepath.Join("screens", file.Name()))
		if err != nil {
			err = fmt.Errorf("failed to read file %s: %v", file.Name(), err)
			return result, err
		}

		// Print the file name and content
		// fmt.Printf("Filename: %s\n", file.Name())
		// fmt.Printf("Content:\n%s\n\n", string(content))

		// Parse JSON content into a map
		jsonData := make(map[string]interface{}, 0)
		if err := json.Unmarshal(content, &jsonData); err != nil {
			err = fmt.Errorf("failed to parse JSON in file %s: %w", file.Name(), err)
			return result, err
		}

		// Append filename to the data
		jsonData["filename"] = file.Name()
		result = append(result, jsonData)

	}

	return
}
