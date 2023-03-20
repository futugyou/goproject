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
	Model            string   `json:"model"`
	Prompt           string   `json:"prompt"`
	Temperature      int32    `json:"temperature"`
	MaxTokens        int32    `json:"max_tokens"`
	Top_p            int32    `json:"top_p"`
	FrequencyPenalty float32  `json:"frequency_penalty"`
	PresencePenalty  float32  `json:"presence_penalty"`
	Stop             []string `json:"stop"`
	BestOf           int      `json:"best_of"`
}

type CreateCompletionRequest struct {
	CompletionModel
}

func (s *CompletionService) CreateCompletion(request CreateCompletionRequest) *openai.CreateCompletionResponse {
	openaikey, _ := config.String("openaikey")
	client := openai.NewClient(openaikey)
	req := openai.CreateCompletionRequest{}
	mapper.AutoMapper(&request, &req)
	result := client.CreateCompletion(req)
	return result
}

func (s *CompletionService) GetExampleSettings(exampleName string) CompletionModel {
	result := CompletionModel{}
	if settings, err := os.ReadFile(fmt.Sprintf("./examples/%s.json", exampleName)); err == nil {

		json.Unmarshal(settings, &result)
	}

	return result
}
