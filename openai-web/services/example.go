package services

import (
	"encoding/json"
	"fmt"
	"os"
)

type ExampleDetailModel struct {
	Model            string   `json:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	MaxTokens        int32    `json:"max_tokens,omitempty"`
	Top_p            float32  `json:"top_p,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
}

type ExampleService struct {
}

func (s *ExampleService) GetExampleSettings(settingName string) ExampleDetailModel {
	if len(settingName) == 0 {
		settingName = "default"
	}

	result := ExampleDetailModel{}
	if settings, err := os.ReadFile(fmt.Sprintf("./examples/%s.json", settingName)); err == nil {
		json.Unmarshal(settings, &result)
	}

	return result
}
