package controllers

import (
	"encoding/json"
	"openai/lib"
	"openai/services"

	"github.com/beego/beego/v2/server/web"
)

// Operations about Chat
type AudioController struct {
	web.Controller
}

// @Title CreateAudioTranscription
// @Description create audio transcription
// @Param	body		body 	lib.CreateAudioTranscriptionRequest	true		"body for create audio transcription content"
// @Success 200 {object} 	lib.CreateAudioTranscriptionResponse
// @router /transcription [post]
func (c *AudioController) CreateAudioTranscription(request lib.CreateAudioTranscriptionRequest) {
	chatService := services.AudioService{}
	var audio lib.CreateAudioTranscriptionRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &audio)
	result := chatService.CreateAudioTranscription(audio)
	c.Ctx.JSONResp(result)
}

// @Title CreateAudioTranslation
// @Description create audio translation
// @Param	body		body 	lib.CreateAudioTranslationRequest	true		"body for create audio translation content"
// @Success 200 {object} 	lib.CreateAudioTranslationResponse
// @router /translation [post]
func (c *AudioController) CreateAudioTranslation(request lib.CreateAudioTranslationRequest) {
	chatService := services.AudioService{}
	var audio lib.CreateAudioTranslationRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &audio)
	result := chatService.CreateAudioTranslation(audio)
	c.Ctx.JSONResp(result)
}
