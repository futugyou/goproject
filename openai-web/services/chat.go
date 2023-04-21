package services

import (
	"strings"
	"time"

	"github.com/beego/beego/v2/core/config"
	lib "github.com/futugyousuzu/go-openai"
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
}

func (s *ChatService) CreateChatCompletion(request lib.CreateChatCompletionRequest) *lib.CreateChatCompletionResponse {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)
	response := client.CreateChatCompletion(request)
	return response
}

func (s *ChatService) CreateChatSSE(request CreateChatRequest) <-chan CreateChatResponse {
	openaikey, _ := config.String("openaikey")
	client := lib.NewClient(openaikey)

	messages := make([]lib.ChatCompletionMessage, 0)
	for i := 0; i < len(request.Messages); i++ {
		switch strings.ToLower(request.Messages[i].Role) {
		case "user":
			messages = append(messages, lib.ChatCompletionMessageFromUser(request.Messages[i].Content))
		case "system":
			messages = append(messages, lib.ChatCompletionMessageFromSystem(request.Messages[i].Content))
		case "assistant":
			messages = append(messages, lib.ChatCompletionMessageFromAssistant(request.Messages[i].Content))
		}
	}

	chatRequest := lib.CreateChatCompletionRequest{
		Model:            request.Model,
		Messages:         messages,
		Temperature:      request.Temperature,
		Top_p:            request.Top_p,
		MaxTokens:        request.MaxTokens,
		PresencePenalty:  request.PresencePenalty,
		FrequencyPenalty: request.FrequencyPenalty,
	}

	stream, err := client.CreateChatStreamCompletion(chatRequest)
	result := make(chan CreateChatResponse)

	if err != nil {
		go func() {
			defer close(result)
			result <- CreateChatResponse{ErrorMessage: err.Error()}
		}()

		return result
	}

	go func() {
		defer close(result)
		defer stream.Close()

		for {
			if !stream.CanReadStream() {
				break
			}

			response := &lib.CreateChatCompletionResponse{}
			ch := CreateChatResponse{}

			if err = stream.ReadStream(response); err != nil {
				ch.ErrorMessage = err.Error()
			} else {
				if response.Created != 0 {
					ch.Created = time.Unix((int64)(response.Created), 0).Format(time.DateTime)
				}

				if response.Usage != nil {
					ch.TotalTokens = response.Usage.TotalTokens
					ch.CompletionTokens = response.Usage.CompletionTokens
					ch.PromptTokens = response.Usage.PromptTokens
				}

				if response.Choices != nil {
					messages := make([]Chat, 0)
					for i := 0; i < len(response.Choices); i++ {
						for j := 0; j < len(response.Choices[i].Message); j++ {
							message := Chat{
								Role:    string(response.Choices[i].Message[j].Role),
								Content: response.Choices[i].Message[j].Content,
							}
							messages = append(messages, message)
						}
					}

					ch.Messages = messages
				}
			}

			result <- ch
		}
	}()

	return result
}
