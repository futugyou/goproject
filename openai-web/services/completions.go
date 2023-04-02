package services

import (
	"encoding/json"
	"fmt"
	"os"

	openai "github.com/futugyousuzu/go-openai"

	"github.com/beego/beego/v2/core/config"
	"github.com/devfeel/mapper"
)

type CompletionService struct {
}

type CompletionModel struct {
	Model            string   `json:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	MaxTokens        int32    `json:"max_tokens,omitempty"`
	Top_p            float32  `json:"top_p,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	BestOf           int32    `json:"best_of,omitempty"`
}

type CreateCompletionRequest struct {
	CompletionModel
}

func (s *CompletionService) CreateCompletion(request CreateCompletionRequest) *openai.CreateCompletionResponse {
	openaikey, _ := config.String("openaikey")
	client := openai.NewClient(openaikey)
	req := openai.CreateCompletionRequest{}
	mapper.AutoMapper(&request.CompletionModel, &req)
	result := client.CreateCompletion(req)
	return result
}

func (s *CompletionService) GetExampleSettings(settingName string) CompletionModel {
	result := CompletionModel{}
	if settings, err := os.ReadFile(fmt.Sprintf("./examples/%s.json", settingName)); err == nil {
		json.Unmarshal(settings, &result)
	}

	return result
}
