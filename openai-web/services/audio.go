package services

import (
	"context"

	openai "github.com/openai/openai-go"
)

type AudioService struct {
	client *openai.Client
}

func NewAudioService(client *openai.Client) *AudioService {
	return &AudioService{
		client: client,
	}
}

func (s *AudioService) CreateAudioTranscription(ctx context.Context, request openai.AudioTranscriptionNewParams) (*openai.Transcription, error) {
	return s.client.Audio.Transcriptions.New(ctx, request)
}

func (s *AudioService) CreateAudioTranslation(ctx context.Context, request openai.AudioTranslationNewParams) (*openai.Translation, error) {
	return s.client.Audio.Translations.New(ctx, request)
}
