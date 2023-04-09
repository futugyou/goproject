package services

import (
	"encoding/json"
	"os"
)

type ExampleModel struct {
	Key              string   `json:"key,omitempty"`
	Title            string   `json:"title,omitempty"`
	SubTitle         string   `json:"subTitle,omitempty"`
	Model            string   `json:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	MaxTokens        int32    `json:"max_tokens,omitempty"`
	Top_p            float32  `json:"top_p,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	Description      string   `json:"description,omitempty"`
	SampleResponse   string   `json:"sample_response,omitempty"`
}

type ExampleService struct {
}

func (s *ExampleService) GetExampleSettings() []ExampleModel {

	result := make([]ExampleModel, 0)
	if settings, err := os.ReadFile("./examples/examples.json"); err == nil {
		json.Unmarshal(settings, &result)
	}

	return result
}
