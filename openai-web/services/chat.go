package services

import (
	"context"
	"time"

	"github.com/openai/openai-go/v3/packages/param"
	openai "github.com/openai/openai-go/v3"
)

type CreateChatRequest struct {
	Model            string  `json:"model"`
	Messages         []Chat  `json:"messages"`
	Temperature      float64 `json:"temperature,omitempty"`
	Top_p            float64 `json:"top_p,omitempty"`
	MaxTokens        int64   `json:"max_tokens,omitempty"`
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"`
}

type CreateChatResponse struct {
	ErrorMessage     string `json:"error,omitempty"`
	Created          string `json:"created,omitempty"`
	PromptTokens     int64  `json:"prompt_tokens,omitempty"`
	CompletionTokens int64  `json:"completion_tokens,omitempty"`
	TotalTokens      int64  `json:"total_tokens,omitempty"`
	Messages         []Chat `json:"messages,omitempty"`
}

type Chat struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatService struct {
	client *openai.Client
}

func NewChatService(client *openai.Client) *ChatService {
	return &ChatService{
		client: client,
	}
}

/*
Request:

	{
		"presence_penalty": 0,
		"prompt": "Say this is a test",
		"model": "gpt-4o",
		"suffix": "string",
		"temperature": 0,
		"top_p": 0,
		"best_of": 1,
		"messages": [{
			"content": "Write me a haiku"
		}],
		"max_tokens": 4096
	}

Response:

	{
	  "created": "2025-03-13 17:26:06",
	  "prompt_tokens": 12,
	  "completion_tokens": 20,
	  "total_tokens": 32,
	  "messages": [
	    {
	      "role": "assistant",
	      "content": "Golden leaves cascade,  \nwhispers of autumn's embrace,  \ntime drifts with the breeze."
	    }
	  ]
	}
*/
func (s *ChatService) CreateChatCompletion(ctx context.Context, req CreateChatRequest) (*CreateChatResponse, error) {
	msg := []openai.ChatCompletionMessageParamUnion{}
	for _, v := range req.Messages {
		if v.Role == "system" {
			msg = append(msg, openai.SystemMessage(v.Content))
		} else {
			msg = append(msg, openai.UserMessage(v.Content))
		}
	}

	request := openai.ChatCompletionNewParams{
		Model:            req.Model,
		Messages:         msg,
		Temperature:      param.NewOpt(req.Temperature),
		TopP:             param.NewOpt(req.Top_p),
		MaxTokens:        param.NewOpt(req.MaxTokens),
		PresencePenalty:  param.NewOpt(req.PresencePenalty),
		FrequencyPenalty: param.NewOpt(req.FrequencyPenalty),
		Seed:             openai.Int(1),
	}
	response, err := s.client.Chat.Completions.New(ctx, request)
	if err != nil {
		return nil, err
	}

	result := &CreateChatResponse{}
	if response.Created != 0 {
		result.Created = time.Unix((int64)(response.Created), 0).Format("2006-01-02 15:04:05")
	}

	result.TotalTokens = response.Usage.TotalTokens
	result.CompletionTokens = response.Usage.CompletionTokens
	result.PromptTokens = response.Usage.PromptTokens

	if response.Choices != nil {
		messages := make([]Chat, 0)
		for i := 0; i < len(response.Choices); i++ {
			message := Chat{
				Role:    string(response.Choices[i].Message.Role),
				Content: response.Choices[i].Message.Content,
			}
			messages = append(messages, message)
		}

		result.Messages = messages
	}
	return result, nil
}

func (s *ChatService) CreateChatSSE(ctx context.Context, req CreateChatRequest) <-chan CreateChatResponse {
	msg := []openai.ChatCompletionMessageParamUnion{}
	for _, v := range req.Messages {
		if v.Role == "system" {
			msg = append(msg, openai.SystemMessage(v.Content))
		} else {
			msg = append(msg, openai.UserMessage(v.Content))
		}
	}

	request := openai.ChatCompletionNewParams{
		Model:            req.Model,
		Messages:         msg,
		Temperature:      param.NewOpt(req.Temperature),
		TopP:             param.NewOpt(req.Top_p),
		MaxTokens:        param.NewOpt(req.MaxTokens),
		PresencePenalty:  param.NewOpt(req.PresencePenalty),
		FrequencyPenalty: param.NewOpt(req.FrequencyPenalty),
	}
	stream := s.client.Chat.Completions.NewStreaming(ctx, request)

	result := make(chan CreateChatResponse)

	go func() {
		defer close(result)
		defer stream.Close()
		for stream.Next() {
			response := stream.Current()
			ch := CreateChatResponse{}
			if response.Created != 0 {
				ch.Created = time.Unix((int64)(response.Created), 0).Format("2006-01-02 15:04:05")
			}

			ch.TotalTokens = response.Usage.TotalTokens
			ch.CompletionTokens = response.Usage.CompletionTokens
			ch.PromptTokens = response.Usage.PromptTokens

			if response.Choices != nil {
				messages := make([]Chat, 0)
				for i := 0; i < len(response.Choices); i++ {
					message := Chat{
						Role:    string(response.Choices[i].Delta.Role),
						Content: response.Choices[i].Delta.Content,
					}
					messages = append(messages, message)
				}

				ch.Messages = messages
			}
			select {
			case result <- ch:
			case <-ctx.Done():
				result <- CreateChatResponse{ErrorMessage: ctx.Err().Error()}
				return
			}
		}

		if err := stream.Err(); err != nil {
			result <- CreateChatResponse{ErrorMessage: err.Error()}
			return
		}

	}()

	return result
}
