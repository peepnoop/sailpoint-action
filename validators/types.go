package validators

import (
    "encoding/json"
    "bufio"
    "os"
	"fmt"
)

// ValidationRule defines the behavior for validating a specific type
type ValidationRule interface {
    Validate(payload string) (interface{}, error)
}

// BaseValidator contains common validation logic to check if the json file is valid
type BaseValidator struct{}

func (bv *BaseValidator) getDecoder(payload string) (*json.Decoder, error) {
    if payload != "" {
        file, err := os.Open(payload)
        if err != nil {
            return nil, err
        }
        defer file.Close()
        return json.NewDecoder(bufio.NewReader(file)), nil
    }
    return json.NewDecoder(bufio.NewReader(os.Stdin)), nil
}

// ValidatorFactory creates appropriate validator based on action
type ValidatorFactory struct {
    validators map[string]ValidationRule
}

// When adding new validator files be sure to add them to the factory map
func NewValidatorFactory() *ValidatorFactory {
    return &ValidatorFactory{
        validators: map[string]ValidationRule{
            "create-transform": &TransformValidator{},
            "update-transform": &TransformUpdateValidator{},
            "create-role":     &RoleValidator{},
        },
    }
}

func (vf *ValidatorFactory) Validate(action string, payload string) (interface{}, error) {
    validator, exists := vf.validators[action]
    if !exists {
        return nil, fmt.Errorf("unknown action: %s", action)
    }
	return validator.Validate(payload)
}