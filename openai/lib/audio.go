package lib

import "os"

const audioTranscriptionPath string = "audio/transcriptions"

type CreateAudioTranscriptionRequest struct {
	File           *os.File `json:"file"`
	Model          string   `json:"model"`
	Prompt         string   `json:"prompt,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"` //  json, text, srt, verbose_json, or vtt.
	Temperature    float32  `json:"temperature,omitempty"`
	Language       string   `json:"language,omitempty"`
}

type CreateAudioTranscriptionResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	Text  string       `json:"text,omitempty"`
}

func (client *openaiClient) CreateAudioTranscription(request CreateAudioTranscriptionRequest) *CreateAudioTranscriptionResponse {
	result := &CreateAudioTranscriptionResponse{}
	client.PostWithFile(audioTranscriptionPath, &request, result)
	return result
}
