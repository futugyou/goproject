package lib

import (
	"golang.org/x/exp/slices"
)

const chatCompletionPath string = "chat/completions"

var supportedChatModel = []string{GPT35_turbo, GPT35_turbo_0301}

type ChatRole string

const ChatRoleSystem ChatRole = "system"
const ChatRoleUser ChatRole = "user"
const ChatRoleAssistant ChatRole = "assistant"

var supportedChatRoles = []ChatRole{ChatRoleSystem, ChatRoleUser, ChatRoleAssistant}

type CreateChatCompletionRequest struct {
	Model            string                  `json:"model"`
	Messages         []chatCompletionMessage `json:"messages"`
	Temperature      float32                 `json:"temperature,omitempty"`
	Top_p            float32                 `json:"top_p,omitempty"`
	N                int32                   `json:"n,omitempty"`
	Stream           bool                    `json:"stream,omitempty"`
	Stop             []string                `json:"stop,omitempty"`
	MaxTokens        int32                   `json:"max_tokens,omitempty"`
	PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int32        `json:"logit_bias,omitempty"`
	User             string                  `json:"user,omitempty"`
}

type chatCompletionMessage struct {
	Role    ChatRole `json:"role,omitempty"`
	Content string   `json:"content,omitempty"`
}

func NewChatCompletionMessage(role ChatRole, message string) chatCompletionMessage {
	return chatCompletionMessage{
		Role:    role,
		Content: message,
	}
}

func ChatCompletionMessageFromUser(message string) chatCompletionMessage {
	return chatCompletionMessage{
		Role:    ChatRoleUser,
		Content: message,
	}
}

func ChatCompletionMessageFromSystem(message string) chatCompletionMessage {
	return chatCompletionMessage{
		Role:    ChatRoleSystem,
		Content: message,
	}
}

func ChatCompletionMessageFromAssistant(message string) chatCompletionMessage {
	return chatCompletionMessage{
		Role:    ChatRoleAssistant,
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

func (c *openaiClient) CreateChatCompletion(request CreateChatCompletionRequest) *CreateChatCompletionResponse {
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

	request.Stream = false
	c.httpClient.Post(chatCompletionPath, request, result)
	return result
}

func validateChatModel(model string) *OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedChatModel, model) {
		return UnsupportedTypeError("Model", model, supportedChatModel)
	}

	return nil
}

func validateChatRole(messages []chatCompletionMessage) *OpenaiError {
	if len(messages) == 0 {
		return MessageError("messages can not be nil.")
	}

	for _, message := range messages {
		if !slices.Contains(supportedChatRoles, message.Role) {
			return UnsupportedTypeError("Message role", message.Role, supportedChatRoles)
		}
	}
	return nil
}

func (c *openaiClient) CreateChatStreamCompletion(request CreateChatCompletionRequest) []*CreateChatCompletionResponse {
	result := make([]*CreateChatCompletionResponse, 0)

	err := validateChatModel(request.Model)
	if err != nil {
		result = append(result, &CreateChatCompletionResponse{Error: err})
		return result
	}

	err = validateChatRole(request.Messages)
	if err != nil {
		result = append(result, &CreateChatCompletionResponse{Error: err})
		return result
	}

	request.Stream = true

	c.httpClient.PostStream(chatCompletionPath, request)

	defer c.httpClient.Close()

	for {
		if c.httpClient.StreamEnd {
			break
		}

		response := &CreateChatCompletionResponse{}
		c.httpClient.ReadStream(response)
		if !c.httpClient.StreamEnd {
			result = append(result, response)
		}
	}

	return result
}
