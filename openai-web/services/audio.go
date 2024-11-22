package services

import (
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

func (s *AudioService) CreateAudioTranscription(request openai.CreateAudioTranscriptionRequest) *openai.CreateAudioTranscriptionResponse {
	response := s.client.Audio.CreateAudioTranscription(request)
	return response
}

func (s *AudioService) CreateAudioTranslation(request openai.CreateAudioTranslationRequest) *openai.CreateAudioTranslationResponse {
	response := s.client.Audio.CreateAudioTranslation(request)
	return response
}
