package ai

import (
	"encoding/json"
)

func AsJson[T any](v T) string {
	if any(v) == nil {
		return "{}"
	}

	if result, err := json.Marshal(v); err != nil {
		return "{}"
	} else {
		return string(result)
	}
}
