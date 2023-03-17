package openai

import (
	"os"
	"strings"

	types "github.com/futugyousuzu/go-openai/audioformattype"
	e "github.com/futugyousuzu/go-openai/internal"
	"golang.org/x/exp/slices"
)

const audioTranscriptionPath string = "audio/transcriptions"
const audioTranslationPath string = "audio/translations"

var supportedAudioType = []string{"mp3", "mp4", "mpeg", "mpga", "m4a", "wav", "webm"}
var supportedAudioModel = []string{
	Whisper_1,
}

type CreateAudioTranscriptionRequest struct {
	File           *os.File              `json:"file"`
	Model          string                `json:"model"`
	Prompt         string                `json:"prompt,omitempty"`
	ResponseFormat types.AudioFormatType `json:"response_format,omitempty"` //  json, text, srt, verbose_json, or vtt.
	Temperature    float32               `json:"temperature,omitempty"`
	Language       string                `json:"language,omitempty"`
}

type CreateAudioTranscriptionResponse struct {
	Error    *e.OpenaiError `json:"error,omitempty"`
	Text     string         `json:"text,omitempty"`
	Task     string         `json:"task,omitempty"`
	Language string         `json:"language,omitempty"`
	Duration float64        `json:"duration,omitempty"`
	Segments []Segments     `json:"segments,omitempty"`
}

type Segments struct {
	ID               int     `json:"id"`
	Seek             int     `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int   `json:"tokens"`
	Temperature      float64 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
	Transient        bool    `json:"transient"`
}

type CreateAudioTranslationRequest struct {
	File           *os.File              `json:"file"`
	Model          string                `json:"model"`
	Prompt         string                `json:"prompt,omitempty"`
	ResponseFormat types.AudioFormatType `json:"response_format,omitempty"` //  json, text, srt, verbose_json, or vtt.
	Temperature    float32               `json:"temperature,omitempty"`
}

type CreateAudioTranslationResponse struct {
	Error *e.OpenaiError `json:"error,omitempty"`
	Text  string         `json:"text,omitempty"`
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

	if request.ResponseFormat == "verbose_json" || request.ResponseFormat == "json" {
		c.httpClient.PostWithFile(audioTranscriptionPath, &request, result)
	} else {
		if err := c.httpClient.PostWithFile(audioTranscriptionPath, &request, &result.Text); err != nil {
			result.Error = e.SystemError(err.Error())
		}
	}

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

	if request.ResponseFormat == "verbose_json" || request.ResponseFormat == "json" {
		c.httpClient.PostWithFile(audioTranslationPath, &request, result)
	} else {
		if err := c.httpClient.PostWithFile(audioTranslationPath, &request, &result.Text); err != nil {
			result.Error = e.SystemError(err.Error())
		}
	}

	return result
}

func validateAudioResponseFormat(responseFormat types.AudioFormatType) *e.OpenaiError {
	if len(responseFormat) > 0 && !slices.Contains(types.SupportededResponseFormatType, responseFormat) {
		return e.UnsupportedTypeError("ResponseFormat", responseFormat, types.SupportededResponseFormatType)
	}

	return nil
}

func validateAudioModel(model string) *e.OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedAudioModel, model) {
		return e.UnsupportedTypeError("Model", model, supportedAudioModel)
	}

	return nil
}

func validateAudioFile(file os.File) *e.OpenaiError {
	segmentations := strings.Split(file.Name(), ".")
	if len(segmentations) <= 1 {
		return e.UnsupportedTypeError("audio type", "nil", supportedAudioType)
	}

	suffix := strings.ToLower(strings.Split(file.Name(), ".")[len(segmentations)-1])
	if !slices.Contains(supportedAudioType, suffix) {
		return e.UnsupportedTypeError("audio type", suffix, supportedAudioType)
	}

	return nil
}
