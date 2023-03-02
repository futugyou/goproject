package lib

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

const audioTranscriptionPath string = "audio/transcriptions"
const audioTranslationPath string = "audio/translations"

var responseFormatType = []string{"json", "text", "srt", "verbose_json", "vtt"}
var audioType = []string{"mp3", "mp4", "mpeg", "mpga", "m4a", "wav", "webm"}

var responseFormatTypeError = func(message string) *OpenaiError {
	return &OpenaiError{
		ErrorMessage: "response format only json, text, srt, verbose_json, or vtt",
		ErrorType:    "invalid parameters",
		Param:        fmt.Sprintf("current response format is: %s", message),
	}
}

var audioTypeError = func(message string) *OpenaiError {
	return &OpenaiError{
		ErrorMessage: "audio type only mp3, mp4, mpeg, mpga, m4a, wav, or webm",
		ErrorType:    "invalid parameters",
		Param:        fmt.Sprintf("current audio type is: %s", message),
	}
}

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

type CreateAudioTranslationRequest struct {
	File           *os.File `json:"file"`
	Model          string   `json:"model"`
	Prompt         string   `json:"prompt,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"` //  json, text, srt, verbose_json, or vtt.
	Temperature    float32  `json:"temperature,omitempty"`
}

type CreateAudioTranslationResponse struct {
	Error *OpenaiError `json:"error,omitempty"`
	Text  string       `json:"text,omitempty"`
}

func (client *openaiClient) CreateAudioTranscription(request CreateAudioTranscriptionRequest) *CreateAudioTranscriptionResponse {
	result := &CreateAudioTranscriptionResponse{}

	if len(request.ResponseFormat) > 0 && !slices.Contains(responseFormatType, request.ResponseFormat) {
		result.Error = responseFormatTypeError(request.ResponseFormat)
		return result
	}

	segmentations := strings.Split(request.File.Name(), ".")
	if len(segmentations) <= 1 {
		result.Error = audioTypeError("no file extension")
		return result
	}

	suffix := strings.Split(request.File.Name(), ".")[len(segmentations)-1]
	if !slices.Contains(audioType, suffix) {
		result.Error = responseFormatTypeError(suffix)
	}

	client.PostWithFile(audioTranscriptionPath, &request, result)
	return result
}

func (client *openaiClient) CreateAudioTranslation(request CreateAudioTranslationRequest) *CreateAudioTranslationResponse {
	result := &CreateAudioTranslationResponse{}
	client.PostWithFile(audioTranslationPath, &request, result)
	return result
}
