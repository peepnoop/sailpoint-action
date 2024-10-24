package validators

import (
    "fmt"
	
	"github.com/sailpoint-oss/golang-sdk/v2/api_v2024"
)

type RoleValidator struct {
    BaseValidator
}

func (rv *RoleValidator) Validate(payload string) (interface{}, error) {
    var role api_v2024.Role
    decoder, err := rv.getDecoder(payload)
    if err != nil {
        return nil, err
    }

    if err := decoder.Decode(&role); err != nil {
        return nil, err
    }

    if role.GetName() == "" {
        return nil, fmt.Errorf("the role must have a name")
    }
    return role, nil
}