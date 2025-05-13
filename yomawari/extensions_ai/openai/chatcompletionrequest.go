package openai

import (
	"encoding/json"
	"errors"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
)

type OpenAIChatCompletionRequest struct {
	Messages []chatcompletion.ChatMessage `json:"messages"`
	Options  chatcompletion.ChatOptions   `json:"options"`
	Stream   bool                         `json:"stream"`
	ModelId  *string                      `json:"model_id"`
}

// UnmarshalJSON
func (r *OpenAIChatCompletionRequest) UnmarshalJSON(data []byte) error {
	var options chatcompletion.ChatOptions
	if err := json.Unmarshal(data, &options); err != nil {
		return err
	}
	// TODO: *r = ....
	return nil
}

// MarshalJSON
func (r OpenAIChatCompletionRequest) MarshalJSON() ([]byte, error) {
	return nil, errors.New("request body serialization is not supported")
}
