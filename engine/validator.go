package main

import (
	"encoding/json"
	"fmt"
)





func GetJSONField(data []byte, field string) (any, error) {
	var obj map[string]any

	if err := json.Unmarshal(data, &obj); err != nil {
		return nil, fmt.Errorf("Engine: invalid JSON: %w", err)
	}

	value, exists := obj[field]
	if !exists {
		return nil, fmt.Errorf("Engine: field %q not found in object: %v", field, obj)
	}

	return value, nil
}

