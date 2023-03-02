package services

import (
	"openai/lib"

	"github.com/beego/beego/v2/core/config"
)

type ChatService struct {
}

func (s *ChatService) CreateChatCompletion(request lib.CreateChatCompletionRequest) *lib.CreateChatCompletionResponse {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	response := client.CreateChatCompletion(request)
	return response
}
