package openai

import (
	"context"

	"golang.org/x/exp/slices"

	role "github.com/futugyousuzu/go-openai/chatrole"
)

type ChatService service

const chatCompletionPath string = "chat/completions"

var supportedChatModel = []string{
	GPT_4,
	GPT_4_0314,
	GPT_4_32k,
	GPT_4_32k_0314,
	GPT35_turbo,
	GPT35_turbo_0301,
}

type CreateChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []ChatCompletionMessage `json:"messages"`
	Temperature      float32                 `json:"temperature,omitempty"`
	Top_p            float32                 `json:"top_p,omitempty"`
	N                int32                   `json:"n,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	MaxTokens        int32                   `json:"max_tokens,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int32        `json:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type chatCompletionRequest struct {
	CreateChatCompletionRequest
	Stream bool `json:"stream,omitempty"`
}

type ChatCompletionMessage struct {
	Role    role.ChatRole `json:"role,omitempty"`
	Content string        `json:"content,omitempty"`
	Name    string        `json:"name,omitempty"`
}

func NewChatCompletionMessage(role role.ChatRole, message string) ChatCompletionMessage {
	return ChatCompletionMessage{
		Role:    role,
		Content: message,
	}
}

func ChatCompletionMessageFromUser(message string) ChatCompletionMessage {
	return ChatCompletionMessage{
		Role:    role.ChatRoleUser,
		Content: message,
	}
}

func ChatCompletionMessageFromSystem(message string) ChatCompletionMessage {
	return ChatCompletionMessage{
		Role:    role.ChatRoleSystem,
		Content: message,
	}
}

func ChatCompletionMessageFromAssistant(message string) ChatCompletionMessage {
	return ChatCompletionMessage{
		Role:    role.ChatRoleAssistant,
		Content: message,
	}
}

type CreateChatCompletionResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	ID      string       `json:"id,omitempty"`
	Object  string       `json:"object,omitempty"`
	Created int32        `json:"created,omitempty"`
	Model   string       `json:"model,omitempty"`
	Choices []Choices    `json:"choices,omitempty"`
	Usage   *Usage       `json:"usage,omitempty"`
}

func (c *ChatService) CreateChatCompletion(ctx context.Context, request CreateChatCompletionRequest) *CreateChatCompletionResponse {
	result := &CreateChatCompletionResponse{}

	err := validateChatModel(request.Model)
	if err != nil {
		result.Error = err
		return result
	}

	err = validateChatRole(request.Messages)
	if err != nil {
		result.Error = err
		return result
	}

	newRequest := chatCompletionRequest{
		CreateChatCompletionRequest: request,
		Stream:                      false,
	}

	c.client.httpClient.Post(ctx, chatCompletionPath, newRequest, result)
	return result
}

func validateChatModel(model string) *OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedChatModel, model) {
		return unsupportedTypeError("Model", model, supportedChatModel)
	}

	return nil
}

func validateChatRole(messages []ChatCompletionMessage) *OpenaiError {
	if len(messages) == 0 {
		return messageError("messages can not be nil.")
	}

	for _, message := range messages {
		if !slices.Contains(role.SupportedChatRoles, message.Role) {
			return unsupportedTypeError("Message role", message.Role, role.SupportedChatRoles)
		}
	}
	return nil
}

// you can read stream in this way.
//
// stream,err:= openai.CreateChatStreamCompletion(CreateChatCompletionRequest{})
//
//	if err!=nil {
//		doSomething()
//	}
//
// defer stream.Close()
//
// result := make([]*CreateChatCompletionResponse, 0)
//
//	for {
//			if !stream.CanReadStream() {
//				break
//			}
//			response := &CreateChatCompletionResponse{}
//			if err=stream.ReadStream(response);err!=nil {
//				doSomething()
//			}else{
//				result = append(result, response)
//			}
//		}
func (c *ChatService) CreateChatStreamCompletion(ctx context.Context, request CreateChatCompletionRequest) (*StreamResponse, *OpenaiError) {
	err := validateChatModel(request.Model)
	if err != nil {
		return nil, err
	}

	err = validateChatRole(request.Messages)
	if err != nil {
		return nil, err
	}

	newRequest := chatCompletionRequest{
		CreateChatCompletionRequest: request,
		Stream:                      true,
	}

	return c.client.httpClient.PostStream(ctx, chatCompletionPath, newRequest)
}
