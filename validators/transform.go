package validators

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2024"

)

type TransformValidator struct {
    BaseValidator
}

// validateFilePath checks if the file exists and has a .json extension
func (tv *TransformValidator) validateFilePath(path string) error {
	// set log level
	log.SetLevel(log.DebugLevel)

    if path == "" {
		log.Debug("File path is empty",path)
        return fmt.Errorf("file path cannot be empty")
    }

    // Check file extension
    if !strings.HasSuffix(strings.ToLower(path), ".json") {
		log.Debug("File path does not have .json extension", path)
        return fmt.Errorf("file must have .json extension")
    }

    // Check if file exists and is readable
    fileInfo, err := os.Stat(path)
    if err != nil {
        if os.IsNotExist(err) {
			log.Debug("File does not exist", path)
            return fmt.Errorf("file does not exist: %s", path)
        }
		log.Debug("Error accessing file", path, err)
        return fmt.Errorf("error accessing file: %v", err)
    }

    // Check if it's actually a file (not a directory)
    if fileInfo.IsDir() {
		log.Debug("File path points to a directory", path)
        return fmt.Errorf("path points to a directory, not a file: %s", path)
    }

    return nil
}

// checks if the file contains valid JSON and has required fields
func (tv *TransformValidator) validateJSONStructure(r io.Reader) error {
    // Read the entire file content
    content, err := io.ReadAll(r)
    if err != nil {
		log.Debug("Error reading file", err)
        return fmt.Errorf("error reading file: %v", err)
    }

    // First, verify it's a valid JSON
    var rawJSON map[string]interface{}
    if err := json.Unmarshal(content, &rawJSON); err != nil {
		log.Debug("Invalid JSON format", err)
        return fmt.Errorf("invalid JSON format: %v", err)
    }

    // Check for required "name" field
    if _, exists := rawJSON["name"]; !exists {
		log.Debug("JSON does not contain a 'name' field", rawJSON)
        return fmt.Errorf("JSON must contain a 'name' field")
    }

    // Check if name is not empty
    if name, ok := rawJSON["name"].(string); !ok || strings.TrimSpace(name) == "" {
		log.Debug("JSON 'name' field is empty", rawJSON)
        return fmt.Errorf("'name' field must be a non-empty string")
    }

    return nil
}

func (tv *TransformValidator) Validate(payload string) (interface{}, error) {
    // First validate the file path
    if err := tv.validateFilePath(payload); err != nil {
        return nil, err
    }

    // Open the file
    file, err := os.Open(payload)
    if err != nil {
        return nil, fmt.Errorf("error opening file: %v", err)
    }
    defer file.Close()

    // Create a copy of the reader since we need to read it twice
    // Once for structure validation and once for actual parsing
    content, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("error reading file: %v", err)
    }

    // Validate JSON structure
    if err := tv.validateJSONStructure(strings.NewReader(string(content))); err != nil {
        return nil, err
    }

    // Now parse into the actual Transform struct
    var transform api_v2024.Transform
    if err := json.NewDecoder(strings.NewReader(string(content))).Decode(&transform); err != nil {
        return nil, fmt.Errorf("error parsing JSON into Transform struct: %v", err)
    }

    // Additional Transform-specific validation can go here
    if err := tv.validateTransformFields(&transform); err != nil {
        return nil, err
    }

    return transform, nil
}

// validateTransformFields performs Transform-specific validation
func (tv *TransformValidator) validateTransformFields(transform *api_v2024.Transform) error {
    if transform == nil {
        return fmt.Errorf("transform cannot be nil")
    }

    // Validate required fields
    if transform.GetName() == "" {
        return fmt.Errorf("transform name cannot be empty")
    }

    return nil
}

type TransformUpdateValidator struct {
    BaseValidator
}

func (tv *TransformUpdateValidator) Validate(payload string) (interface{}, error) {
    var transform api_v2024.TransformRead
    decoder, err := tv.getDecoder(payload)
    if err != nil {
        return nil, err
    }

    if err := decoder.Decode(&transform); err != nil {
        return nil, err
    }

    if transform.Id == "" {
        return nil, fmt.Errorf("the transform must have an ID")
    }
    return transform.Id, nil
}