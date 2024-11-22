package services

import (
	"context"
	"os"

	openai "github.com/futugyousuzu/go-openai"
)

type AudioService struct {
	client *openai.OpenaiClient
}

func NewAudioService(client *openai.OpenaiClient) *AudioService {
	if client == nil {
		openaikey := os.Getenv("openaikey")
		client = openai.NewClient(openaikey)
	}
	return &AudioService{
		client: client,
	}
}

func (s *AudioService) CreateAudioTranscription(ctx context.Context, request openai.CreateAudioTranscriptionRequest) *openai.CreateAudioTranscriptionResponse {
	response := s.client.Audio.CreateAudioTranscription(ctx, request)
	return response
}

func (s *AudioService) CreateAudioTranslation(ctx context.Context, request openai.CreateAudioTranslationRequest) *openai.CreateAudioTranslationResponse {
	response := s.client.Audio.CreateAudioTranslation(ctx, request)
	return response
}
