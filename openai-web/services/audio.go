package services

import (
	lib "github.com/futugyousuzu/go-openai"

	"github.com/beego/beego/v2/core/config"
)

type AudioService struct {
}

func (s *AudioService) CreateAudioTranscription(request lib.CreateAudioTranscriptionRequest) *lib.CreateAudioTranscriptionResponse {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	response := client.CreateAudioTranscription(request)
	return response
}

func (s *AudioService) CreateAudioTranslation(request lib.CreateAudioTranslationRequest) *lib.CreateAudioTranslationResponse {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	response := client.CreateAudioTranslation(request)
	return response
}
