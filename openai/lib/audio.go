package lib

import (
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

const audioTranscriptionPath string = "audio/transcriptions"
const audioTranslationPath string = "audio/translations"

var supportededResponseFormatType = []string{"json", "text", "srt", "verbose_json", "vtt"}
var supportedAudioType = []string{"mp3", "mp4", "mpeg", "mpga", "m4a", "wav", "webm"}
var supportedAudioModel = []string{Whisper_1}

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

func (c *openaiClient) CreateAudioTranscription(request CreateAudioTranscriptionRequest) *CreateAudioTranscriptionResponse {
	result := &CreateAudioTranscriptionResponse{}

	err := validateAudioModel(request.Model)
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

	c.httpClient.PostWithFile(audioTranscriptionPath, &request, result)
	return result
}

func (c *openaiClient) CreateAudioTranslation(request CreateAudioTranslationRequest) *CreateAudioTranslationResponse {
	result := &CreateAudioTranslationResponse{}

	err := validateAudioModel(request.Model)
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

	c.httpClient.PostWithFile(audioTranslationPath, &request, result)
	return result
}

func validateAudioResponseFormat(responseFormat string) *OpenaiError {
	if len(responseFormat) > 0 && !slices.Contains(supportededResponseFormatType, responseFormat) {
		return UnsupportedTypeError("ResponseFormat", responseFormat, supportededResponseFormatType)
	}

	return nil
}

func validateAudioModel(model string) *OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedAudioModel, model) {
		return UnsupportedTypeError("Model", model, supportedAudioModel)
	}

	return nil
}

func validateAudioFile(file os.File) *OpenaiError {
	segmentations := strings.Split(file.Name(), ".")
	if len(segmentations) <= 1 {
		return UnsupportedTypeError("audio type", "nil", supportedAudioType)
	}

	suffix := strings.ToLower(strings.Split(file.Name(), ".")[len(segmentations)-1])
	if !slices.Contains(supportedAudioType, suffix) {
		return UnsupportedTypeError("audio type", suffix, supportedAudioType)
	}

	return nil
}
