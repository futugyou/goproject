package services

import (
	"context"
	"time"

	openai "github.com/openai/openai-go"
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

func (s *ChatService) CreateChatCompletion(ctx context.Context, request openai.ChatCompletionNewParams) (*openai.ChatCompletion, error) {
	return s.client.Chat.Completions.New(ctx, request)
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

	messages := openai.F(msg)
	request := openai.ChatCompletionNewParams{
		Model:            openai.F(req.Model),
		Messages:         messages,
		Temperature:      openai.F(req.Temperature),
		TopP:             openai.F(req.Top_p),
		MaxTokens:        openai.F(req.MaxTokens),
		PresencePenalty:  openai.F(req.PresencePenalty),
		FrequencyPenalty: openai.F(req.FrequencyPenalty),
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
			close(result)
			return
		}

	}()

	return result
}
