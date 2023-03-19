package services

import (
	openai "github.com/futugyousuzu/go-openai"

	"github.com/beego/beego/v2/core/config"
	"github.com/devfeel/mapper"
)

type CompletionService struct {
}

type CreateCompletionRequest struct {
	Prompt           string  `json:"prompt"`
	MaxTokens        int     `json:"max_tokens"`
	Temperature      float64 `json:"temperature"`
	TopP             int     `json:"top_p"`
	FrequencyPenalty int     `json:"frequency_penalty"`
	PresencePenalty  int     `json:"presence_penalty"`
	BestOf           int     `json:"best_of"`
}

func (s *CompletionService) CreateCompletion(request CreateCompletionRequest) *openai.CreateCompletionResponse {
	openaikey, _ := config.String("openaikey")
	client := openai.NewClient(openaikey)
	req := openai.CreateCompletionRequest{}
	mapper.AutoMapper(&request, &req)
	result := client.CreateCompletion(req)
	return result
}
