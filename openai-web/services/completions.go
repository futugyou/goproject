package services

import (
	"context"
	"os"
	"time"

	openai "github.com/futugyou/ai-extension/openai"

	"github.com/devfeel/mapper"
)

type CompletionService struct {
	client *openai.OpenaiClient
}

func NewCompletionService(client *openai.OpenaiClient) *CompletionService {
	if client == nil {
		openaikey := os.Getenv("openaikey")
		client = openai.NewClient(openaikey)
	}
	return &CompletionService{
		client: client,
	}
}

type CompletionModel struct {
	Model            string   `json:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	MaxTokens        int32    `json:"max_tokens,omitempty"`
	Top_p            float32  `json:"top_p,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	BestOf           int32    `json:"best_of,omitempty"`
	Suffix           string   `json:"suffix,omitempty"`
}

type CreateCompletionRequest struct {
	CompletionModel
}

type CreateCompletionResponse struct {
	ErrorMessage     string   `json:"error,omitempty"`
	Created          string   `json:"created,omitempty"`
	PromptTokens     int32    `json:"prompt_tokens,omitempty"`
	CompletionTokens int32    `json:"completion_tokens,omitempty"`
	TotalTokens      int32    `json:"total_tokens,omitempty"`
	Texts            []string `json:"texts,omitempty"`
}

func (s *CompletionService) CreateCompletion(ctx context.Context, request CreateCompletionRequest) CreateCompletionResponse {
	req := openai.CreateCompletionRequest{}
	mapper.AutoMapper(&request.CompletionModel, &req)
	apiresult := s.client.Completion.CreateCompletion(ctx, req)

	result := CreateCompletionResponse{}
	if apiresult != nil {
		if apiresult.Error != nil {
			result.ErrorMessage = apiresult.Error.Error()
		}

		if apiresult.Created != 0 {
			result.Created = time.Unix((int64)(apiresult.Created), 0).Format("2006-01-02 15:04:05")
		}

		if apiresult.Usage != nil {
			result.TotalTokens = apiresult.Usage.TotalTokens
			result.CompletionTokens = apiresult.Usage.CompletionTokens
			result.PromptTokens = apiresult.Usage.PromptTokens
		}

		if apiresult.Choices != nil {
			texts := make([]string, 0)
			for i := 0; i < len(apiresult.Choices); i++ {
				texts = append(texts, apiresult.Choices[i].Text)
			}

			result.Texts = texts
		}
	}

	return result
}
func (s *CompletionService) CreateCompletionSSE(ctx context.Context, request CreateCompletionRequest) <-chan CreateCompletionResponse {
	req := openai.CreateCompletionRequest{}
	mapper.AutoMapper(&request.CompletionModel, &req)

	result := make(chan CreateCompletionResponse)
	stream, err := s.client.Completion.CreateStreamCompletion(ctx, req)
	if err != nil {
		result <- CreateCompletionResponse{ErrorMessage: err.Error()}
		close(result)
		return result
	}

	go func() {
		defer close(result)
		defer stream.Close()

		for {
			select {
			case <-ctx.Done():
				result <- CreateCompletionResponse{ErrorMessage: ctx.Err().Error()}
				return
			default:
			}

			if !stream.CanReadStream() {
				break
			}

			response := &openai.CreateCompletionResponse{}
			ch := CreateCompletionResponse{}

			if err := stream.ReadStream(ctx, response); err != nil {
				ch.ErrorMessage = err.Error()
			} else {
				if response.Created != 0 {
					ch.Created = time.Unix(int64(response.Created), 0).Format("2006-01-02 15:04:05")
				}

				if response.Usage != nil {
					ch.TotalTokens = response.Usage.TotalTokens
					ch.CompletionTokens = response.Usage.CompletionTokens
					ch.PromptTokens = response.Usage.PromptTokens
				}

				if response.Choices != nil {
					texts := make([]string, len(response.Choices))
					for i, choice := range response.Choices {
						texts[i] = choice.Text
					}
					ch.Texts = texts
				}
			}

			select {
			case result <- ch:
			case <-ctx.Done():
				result <- CreateCompletionResponse{ErrorMessage: ctx.Err().Error()}
				return
			}
		}
	}()

	return result
}
