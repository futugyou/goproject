package services

import (
	"context"
	"os"

	openai "github.com/futugyou/ai-extension/openai"
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

func (s *FineTuneService) ListFinetuneEvents(ctx context.Context, fine_tune_id string) *openai.ListFinetuneEventResponse {
	result := s.client.Finetune.ListFinetuneEvents(ctx, fine_tune_id)
	return result
}
