package lib

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

const audioTranscriptionPath string = "audio/transcriptions"
const audioTranslationPath string = "audio/translations"

var supportededResponseFormatType = []string{"json", "text", "srt", "verbose_json", "vtt"}
var supportedAudioType = []string{"mp3", "mp4", "mpeg", "mpga", "m4a", "wav", "webm"}
var supportedAudioModel = []string{"whisper-1"}

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

var audioModelError = func(message string) *OpenaiError {
	return &OpenaiError{
		ErrorMessage: "Only whisper-1 is currently available.",
		ErrorType:    "invalid parameters",
		Param:        fmt.Sprintf("current model is: %s", message),
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

	err := validateAudioModel(request.ResponseFormat)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateAudioResponseFormat(request.ResponseFormat)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateAudioFile(*request.File)
	if err != nil {
		result.Error = err
		return result
	}

	client.PostWithFile(audioTranscriptionPath, &request, result)
	return result
}

func (client *openaiClient) CreateAudioTranslation(request CreateAudioTranslationRequest) *CreateAudioTranslationResponse {
	result := &CreateAudioTranslationResponse{}

	err := validateAudioModel(request.ResponseFormat)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateAudioResponseFormat(request.ResponseFormat)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateAudioFile(*request.File)
	if err != nil {
		result.Error = err
		return result
	}

	client.PostWithFile(audioTranslationPath, &request, result)
	return result
}

func validateAudioResponseFormat(responseFormat string) *OpenaiError {
	if len(responseFormat) > 0 && !slices.Contains(supportededResponseFormatType, responseFormat) {
		return responseFormatTypeError(responseFormat)
	}

	return nil
}

func validateAudioModel(model string) *OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedAudioModel, model) {
		return audioModelError(model)
	}

	return nil
}

func validateAudioFile(file os.File) *OpenaiError {
	segmentations := strings.Split(file.Name(), ".")
	if len(segmentations) <= 1 {
		return audioTypeError("no file extension")
	}

	suffix := strings.ToLower(strings.Split(file.Name(), ".")[len(segmentations)-1])
	if !slices.Contains(supportedAudioType, suffix) {
		return responseFormatTypeError(suffix)
	}

	return nil
}
