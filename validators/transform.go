package validators

import (
    "fmt"

	"github.com/sailpoint-oss/golang-sdk/v2/api_v2024"
)

type TransformValidator struct {
    BaseValidator
}

func (tv *TransformValidator) Validate(payload string) (interface{}, error) {
    var transform api_v2024.Transform
    decoder, err := tv.getDecoder(payload)
    if err != nil {
        return nil, err
    }

    if err := decoder.Decode(&transform); err != nil {
        return nil, err
    }

    if transform.GetName() == "" {
        return nil, fmt.Errorf("the transform must have a name")
    }
    return transform, nil
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