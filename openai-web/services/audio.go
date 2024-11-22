package services

import (
	"os"

	lib "github.com/futugyousuzu/go-openai"
)

type AudioService struct {
}

func (s *AudioService) CreateAudioTranscription(request lib.CreateAudioTranscriptionRequest) *lib.CreateAudioTranscriptionResponse {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	response := client.Audio.CreateAudioTranscription(request)
	return response
}

func (s *AudioService) CreateAudioTranslation(request lib.CreateAudioTranslationRequest) *lib.CreateAudioTranslationResponse {
	openaikey := os.Getenv("openaikey")
	client := lib.NewClient(openaikey)
	response := client.Audio.CreateAudioTranslation(request)
	return response
}
