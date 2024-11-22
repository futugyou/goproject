package services

import (
	"os"

	openai "github.com/futugyousuzu/go-openai"
)

type FineTuneService struct {
	client *openai.OpenaiClient
}

func NewFineTuneService(client *openai.OpenaiClient) *FineTuneService {
	if client == nil {
		openaikey := os.Getenv("openaikey")
		client = openai.NewClient(openaikey)
	}
	return &FineTuneService{
		client: client,
	}
}

func (s *FineTuneService) ListFinetuneEvents(fine_tune_id string) *openai.ListFinetuneEventResponse {
	result := s.client.Finetune.ListFinetuneEvents(fine_tune_id)
	return result
}
