package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

// SharedClient is a global HttpClient instance used for local access to non-production environments.
var SharedClient = &http.Client{
	Timeout: 0,
}

// TransferNanosecondsTime transfers the nanosecond time in the response to the metadata dictionary
func TransferNanosecondsTime[T any](response T, getNanoseconds func(T) *int64, key string, metadata *map[string]int64) {
	if duration := getNanoseconds(response); duration != nil {
		if *metadata == nil {
			*metadata = make(map[string]int64)
		}
		(*metadata)[key] = *duration
	}
}

// ThrowUnsuccessfulOllamaResponseAsync handles non-2xx HTTP responses
func ThrowUnsuccessfulOllamaResponseAsync(ctx context.Context, response *http.Response) error {
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	}

	// Read the complete response content
	var buf bytes.Buffer
	_, err := io.Copy(&buf, response.Body)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	errorContent := buf.String()

	// Try parsing the JSON error field
	var jsonData map[string]any
	if err := json.Unmarshal([]byte(errorContent), &jsonData); err == nil {
		if errMsg, ok := jsonData["error"].(string); ok {
			errorContent = errMsg
		}
	}

	return errors.New("Ollama error: " + errorContent)
}

func Ptr[T any](v T) *T {
	return &v
}

func StringToTimePtr(s *string, layout string) *time.Time {
	if s == nil {
		return nil
	}

	parsedTime, err := time.Parse(layout, *s)
	if err != nil {
		return nil
	}

	return &parsedTime
}
