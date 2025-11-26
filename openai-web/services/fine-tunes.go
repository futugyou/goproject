package services

import (
	"context"

	openai "github.com/openai/openai-go/v3"
)

type FineTuneService struct {
	client *openai.Client
}

func NewFineTuneService(client *openai.Client) *FineTuneService {
	return &FineTuneService{
		client: client,
	}
}

func (s *FineTuneService) ListFinetuneEvents(ctx context.Context, fine_tune_file_id string) (*openai.FineTuningJob, error) {
	// file, err := s.client.Files.New(ctx, openai.FileNewParams{
	// 	File:    openai.F[io.Reader](data),
	// 	Purpose: openai.F(openai.FilePurposeFineTune),
	// })
	request := openai.FineTuningJobNewParams{
		Model:        openai.FineTuningJobNewParamsModelGPT3_5Turbo,
		TrainingFile: fine_tune_file_id,
	}
	return s.client.FineTuning.Jobs.New(ctx, request)
}
