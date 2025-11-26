package services

import (
	"context"
	"strings"
	"time"

	"github.com/openai/openai-go/v3/packages/param"
	openai "github.com/openai/openai-go/v3"
)

// CompletionService is deprecated.
// Deprecated: This service is no longer to use, use ChatService
type CompletionService struct {
	client *openai.Client
}

func NewCompletionService(client *openai.Client) *CompletionService {
	return &CompletionService{
		client: client,
	}
}

type CompletionModel struct {
	Model            string   `json:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Temperature      float64  `json:"temperature,omitempty"`
	MaxTokens        int64    `json:"max_tokens,omitempty"`
	Top_p            float64  `json:"top_p,omitempty"`
	FrequencyPenalty float64  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float64  `json:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	BestOf           int64    `json:"best_of,omitempty"`
	Suffix           string   `json:"suffix,omitempty"`
}

type CreateCompletionRequest struct {
	CompletionModel
}

type CreateCompletionResponse struct {
	ErrorMessage     string   `json:"error,omitempty"`
	Created          string   `json:"created,omitempty"`
	PromptTokens     int64    `json:"prompt_tokens,omitempty"`
	CompletionTokens int64    `json:"completion_tokens,omitempty"`
	TotalTokens      int64    `json:"total_tokens,omitempty"`
	Texts            []string `json:"texts,omitempty"`
}

func (s *CompletionService) CreateCompletion(ctx context.Context, request CreateCompletionRequest) (*CreateCompletionResponse, error) {
	model := (openai.CompletionNewParamsModel)(request.Model)
	if len(model) == 0 {
		if strings.HasPrefix(request.Model, "babbage") {
			model = openai.CompletionNewParamsModelBabbage002
		} else if strings.HasPrefix(request.Model, "davinci") {
			model = openai.CompletionNewParamsModelDavinci002
		} else if strings.HasPrefix(request.Model, "gpt-3.5") {
			model = openai.CompletionNewParamsModelGPT3_5TurboInstruct
		}
	}

	prompt := openai.CompletionNewParamsPromptUnion{
		OfArrayOfStrings: []string{request.Prompt},
	}
	stop := openai.CompletionNewParamsStopUnion{
		OfStringArray: request.Stop,
	}
	req := openai.CompletionNewParams{
		Model:            model,
		Prompt:           prompt,
		BestOf:           param.NewOpt(request.BestOf),
		FrequencyPenalty: param.NewOpt(request.FrequencyPenalty),
		MaxTokens:        param.NewOpt(request.MaxTokens),
		PresencePenalty:  param.NewOpt(request.PresencePenalty),
		Stop:             stop,
		Suffix:           param.NewOpt(request.Suffix),
		Temperature:      param.NewOpt(request.Temperature),
		TopP:             param.NewOpt(request.Top_p),
		Seed:             openai.Int(1),
	}

	apiresult, err := s.client.Completions.New(ctx, req)
	if err != nil {
		return nil, err
	}

	result := &CreateCompletionResponse{}
	if apiresult != nil {
		if apiresult.Created != 0 {
			result.Created = time.Unix((int64)(apiresult.Created), 0).Format("2006-01-02 15:04:05")
		}

		result.TotalTokens = apiresult.Usage.TotalTokens
		result.CompletionTokens = apiresult.Usage.CompletionTokens
		result.PromptTokens = apiresult.Usage.PromptTokens

		if apiresult.Choices != nil {
			texts := make([]string, 0)
			for i := 0; i < len(apiresult.Choices); i++ {
				texts = append(texts, apiresult.Choices[i].Text)
			}

			result.Texts = texts
		}
	}

	return result, nil
}

func (s *CompletionService) CreateCompletionSSE(ctx context.Context, request CreateCompletionRequest) <-chan CreateCompletionResponse {
	model := openai.CompletionNewParamsModelGPT3_5TurboInstruct
	if strings.HasPrefix(request.Model, "babbage") {
		model = openai.CompletionNewParamsModelBabbage002
	} else if strings.HasPrefix(request.Model, "davinci") {
		model = openai.CompletionNewParamsModelDavinci002
	}

	prompt := openai.CompletionNewParamsPromptUnion{
		OfArrayOfStrings: []string{request.Prompt},
	}
	stop := openai.CompletionNewParamsStopUnion{
		OfStringArray: request.Stop,
	}
	req := openai.CompletionNewParams{
		Model:            model,
		Prompt:           prompt,
		BestOf:           param.NewOpt(request.BestOf),
		FrequencyPenalty: param.NewOpt(request.FrequencyPenalty),
		MaxTokens:        param.NewOpt(request.MaxTokens),
		PresencePenalty:  param.NewOpt(request.PresencePenalty),
		Stop:             stop,
		Suffix:           param.NewOpt(request.Suffix),
		Temperature:      param.NewOpt(request.Temperature),
		TopP:             param.NewOpt(request.Top_p),
	}

	stream := s.client.Completions.NewStreaming(ctx, req)
	result := make(chan CreateCompletionResponse)

	go func() {
		defer close(result)
		defer stream.Close()
		for stream.Next() {
			response := stream.Current()
			ch := CreateCompletionResponse{}
			if response.Created != 0 {
				ch.Created = time.Unix(int64(response.Created), 0).Format("2006-01-02 15:04:05")
			}

			ch.TotalTokens = response.Usage.TotalTokens
			ch.CompletionTokens = response.Usage.CompletionTokens
			ch.PromptTokens = response.Usage.PromptTokens

			if response.Choices != nil {
				texts := make([]string, len(response.Choices))
				for i, choice := range response.Choices {
					texts[i] = choice.Text
				}
				ch.Texts = texts
			}

			select {
			case result <- ch:
			case <-ctx.Done():
				result <- CreateCompletionResponse{ErrorMessage: ctx.Err().Error()}
				return
			}
		}

		if err := stream.Err(); err != nil {
			result <- CreateCompletionResponse{ErrorMessage: err.Error()}
			return
		}

	}()

	return result
}
