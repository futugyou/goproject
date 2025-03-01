package services

import (
	"context"
	"os"
	"strings"
	"time"

	openai "github.com/futugyou/ai-extension/openai"
)

type CreateChatRequest struct {
	Model            string  `json:"model"`
	Messages         []Chat  `json:"messages"`
	Temperature      float32 `json:"temperature,omitempty"`
	Top_p            float32 `json:"top_p,omitempty"`
	MaxTokens        int32   `json:"max_tokens,omitempty"`
	PresencePenalty  float32 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32 `json:"frequency_penalty,omitempty"`
}

type CreateChatResponse struct {
	ErrorMessage     string `json:"error,omitempty"`
	Created          string `json:"created,omitempty"`
	PromptTokens     int32  `json:"prompt_tokens,omitempty"`
	CompletionTokens int32  `json:"completion_tokens,omitempty"`
	TotalTokens      int32  `json:"total_tokens,omitempty"`
	Messages         []Chat `json:"messages,omitempty"`
}

type Chat struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatService struct {
	client *openai.OpenaiClient
}

func NewChatService(client *openai.OpenaiClient) *ChatService {
	if client == nil {
		openaikey := os.Getenv("openaikey")
		client = openai.NewClient(openaikey)
	}
	return &ChatService{
		client: client,
	}
}

func (s *ChatService) CreateChatCompletion(ctx context.Context, request openai.CreateChatCompletionRequest) *openai.CreateChatCompletionResponse {
	response := s.client.Chat.CreateChatCompletion(ctx, request)
	return response
}

func (s *ChatService) CreateChatSSE(ctx context.Context, request CreateChatRequest) <-chan CreateChatResponse {
	messages := make([]openai.ChatCompletionMessage, 0)
	for i := 0; i < len(request.Messages); i++ {
		switch strings.ToLower(request.Messages[i].Role) {
		case "user":
			messages = append(messages, openai.ChatCompletionMessageFromUser(request.Messages[i].Content))
		case "system":
			messages = append(messages, openai.ChatCompletionMessageFromSystem(request.Messages[i].Content))
		case "assistant":
			messages = append(messages, openai.ChatCompletionMessageFromAssistant(request.Messages[i].Content))
		}
	}

	chatRequest := openai.CreateChatCompletionRequest{
		Model:            request.Model,
		Messages:         messages,
		Temperature:      request.Temperature,
		Top_p:            request.Top_p,
		MaxTokens:        request.MaxTokens,
		PresencePenalty:  request.PresencePenalty,
		FrequencyPenalty: request.FrequencyPenalty,
	}

	stream, err := s.client.Chat.CreateChatStreamCompletion(ctx, chatRequest)
	result := make(chan CreateChatResponse)
	if err != nil {
		result <- CreateChatResponse{ErrorMessage: err.Error()}
		close(result)
		return result
	}

	go func() {
		defer close(result)
		defer stream.Close()

		for {
			select {
			case <-ctx.Done():
				result <- CreateChatResponse{ErrorMessage: ctx.Err().Error()}
				return
			default:
			}

			if !stream.CanReadStream() {
				break
			}

			response := &openai.CreateChatCompletionResponse{}
			ch := CreateChatResponse{}

			if err = stream.ReadStream(ctx, response); err != nil {
				ch.ErrorMessage = err.Error()
			} else {
				if response.Created != 0 {
					ch.Created = time.Unix((int64)(response.Created), 0).Format("2006-01-02 15:04:05")
				}

				if response.Usage != nil {
					ch.TotalTokens = response.Usage.TotalTokens
					ch.CompletionTokens = response.Usage.CompletionTokens
					ch.PromptTokens = response.Usage.PromptTokens
				}

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
			}

			select {
			case result <- ch:
			case <-ctx.Done():
				result <- CreateChatResponse{ErrorMessage: ctx.Err().Error()}
				return
			}
		}
	}()

	return result
}
