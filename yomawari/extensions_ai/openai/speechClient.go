package openai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/speechtotext"
	rawopenai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/packages/param"
)

var _ speechtotext.ISpeechToTextClient = (*OpenAISpeechToTextClient)(nil)

type OpenAISpeechToTextClient struct {
	metadata     speechtotext.SpeechToTextClientMetadata
	openAIClient *rawopenai.Client
	modelId      *string
}

func NewOpenAISpeechToTextClient(openAIClient *rawopenai.Client, modelId *string) *OpenAISpeechToTextClient {
	name := "openai"
	return &OpenAISpeechToTextClient{
		metadata: speechtotext.SpeechToTextClientMetadata{
			ProviderName:   &name,
			DefaultModelId: modelId,
		},
		openAIClient: openAIClient,
		modelId:      modelId,
	}
}

// GetStreamingTextWithDataConten implements speechtotext.ISpeechToTextClient.
func (o *OpenAISpeechToTextClient) GetStreamingTextWithDataConten(ctx context.Context, audioSpeechContent contents.DataContent, options *speechtotext.SpeechToTextOptions) (<-chan speechtotext.SpeechToTextResponse, <-chan error) {
	reader := bytes.NewReader(audioSpeechContent.Data)
	return o.GetStreamingText(ctx, io.NopCloser(reader), options)
}

// GetText implements speechtotext.ISpeechToTextClient.
func (o *OpenAISpeechToTextClient) GetText(ctx context.Context, audioSpeechStream io.ReadCloser, options *speechtotext.SpeechToTextOptions) (*speechtotext.SpeechToTextResponse, error) {
	modelid := ""
	if o.modelId != nil {
		modelid = *o.modelId
	}
	if len(modelid) == 0 && options != nil && options.ModelId != nil {
		modelid = *options.ModelId
	}
	prompt := ""
	language := ""
	var responseFormat rawopenai.AudioTranslationNewParamsResponseFormat = "json"
	var responseFormat2 rawopenai.AudioResponseFormat = "json"
	var temperature float64 = 0.0
	var granularities []string = []string{}
	startTime := time.Now().UTC()
	if options != nil && options.AdditionalProperties != nil {
		if v, ok := options.AdditionalProperties["Temperature"].(float64); ok {
			temperature = v
		}
		if v, ok := options.AdditionalProperties["ResponseFormat"].(string); ok {
			//`json`, `text`, `srt`, `verbose_json`, or `vtt`
			switch v {
			case "text":
				responseFormat = "text"
				responseFormat2 = "text"
			case "srt":
				responseFormat = "srt"
				responseFormat2 = "srt"
			case "verbose_json":
				responseFormat = "verbose_json"
				responseFormat2 = "verbose_json"
			case "vtt":
				responseFormat = "vtt"
				responseFormat2 = "vtt"
			default:
				responseFormat = "json"
				responseFormat2 = "json"
			}
		}
		if v, ok := options.AdditionalProperties["Prompt"].(string); ok {
			prompt = v
		}
		if v, ok := options.AdditionalProperties["TimestampGranularities"].(string); ok {
			switch v {
			case "word":
				granularities = []string{"word"}
			case "segment":
				granularities = []string{"segment"}
			}
		}
	}

	text := ""
	if IsTranslationRequest(options) {
		var reader io.Reader = audioSpeechStream
		body := rawopenai.AudioTranscriptionNewParams{
			File:                   reader,
			Model:                  modelid,
			Language:               param.NewOpt(language),
			Prompt:                 param.NewOpt(prompt),
			ResponseFormat:         responseFormat2,
			Temperature:            param.NewOpt(temperature),
			TimestampGranularities: granularities,
		}
		res, err := o.openAIClient.Audio.Transcriptions.New(ctx, body)
		if err != nil {
			return nil, err
		}
		text = res.Text
	} else {
		var reader io.Reader = audioSpeechStream
		body := rawopenai.AudioTranslationNewParams{
			File:           reader,
			Model:          modelid,
			Prompt:         param.NewOpt(prompt),
			ResponseFormat: responseFormat,
			Temperature:    param.NewOpt(temperature),
		}
		res, err := o.openAIClient.Audio.Translations.New(ctx, body)
		if err != nil {
			return nil, err
		}
		text = res.Text
	}

	endTime := time.Now().UTC()
	response := &speechtotext.SpeechToTextResponse{
		StartTime:            &startTime,
		EndTime:              &endTime,
		ModelId:              &modelid,
		Contents:             []contents.IAIContent{contents.NewTextContent(text)},
		AdditionalProperties: map[string]interface{}{"Duration": endTime.Sub(startTime)},
		RawRepresentation:    text,
	}
	if len(language) > 0 {
		response.AdditionalProperties["Language"] = language
	}
	return response, nil
}

// GetTextWithDataContent implements speechtotext.ISpeechToTextClient.
func (o *OpenAISpeechToTextClient) GetTextWithDataContent(ctx context.Context, audioSpeechContent contents.DataContent, options *speechtotext.SpeechToTextOptions) (*speechtotext.SpeechToTextResponse, error) {
	reader := bytes.NewReader(audioSpeechContent.Data)
	return o.GetText(ctx, io.NopCloser(reader), options)
}

func IsTranslationRequest(options *speechtotext.SpeechToTextOptions) bool {
	return options != nil && options.TextLanguage != nil && (options.SpeechLanguage == nil || options.SpeechLanguage != options.TextLanguage)
}

// GetStreamingText implements speechtotext.ISpeechToTextClient.
func (o *OpenAISpeechToTextClient) GetStreamingText(ctx context.Context, audioSpeechStream io.ReadCloser, options *speechtotext.SpeechToTextOptions) (<-chan speechtotext.SpeechToTextResponse, <-chan error) {
	result := make(chan speechtotext.SpeechToTextResponse)
	errCh := make(chan error, 1)
	go func() {
		defer close(result)
		defer close(errCh)
		defer func() {
			if r := recover(); r != nil {
				errCh <- fmt.Errorf("panic recovered: %v", r)
			}
		}()
		resp, err := o.GetText(ctx, audioSpeechStream, options)
		if err != nil {
			errCh <- err
			return
		}
		result <- *resp
	}()
	return result, errCh
}
