package models

import (
	"errors"
	"strings"
)

// CleanName cleans the index name, returning a non-empty value if possible
// name: Input index name
// defaultName: Default value to fall back when input is empty
// Returns: Non-empty index name or error if both inputs are empty
func CleanName(name *string, defaultName *string) (string, error) {
	// Check if both are nil or empty
	if (name == nil || strings.TrimSpace(*name) == "") &&
		(defaultName == nil || strings.TrimSpace(*defaultName) == "") {
		return "", errors.New("both index name and default fallback value are empty. Provide an index name or a default value to use when the index name is empty")
	}

	// Process defaultName
	var cleanDefault string
	if defaultName != nil {
		cleanDefault = strings.TrimSpace(*defaultName)
	}

	// If name is nil, return default
	if name == nil {
		return cleanDefault, nil
	}

	// Process name
	cleanName := strings.TrimSpace(*name)
	if cleanName == "" {
		return cleanDefault, nil
	}

	return cleanName, nil
}
