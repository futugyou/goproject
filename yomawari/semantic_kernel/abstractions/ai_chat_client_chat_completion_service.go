package abstractions

import (
	"context"
	"encoding/json"

	"github.com/futugyou/yomawari/core"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	aicontents "github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

var _ IChatCompletionService = (*ChatClientChatCompletionService)(nil)

type ChatClientChatCompletionService struct {
	chatClient chatcompletion.IChatClient
	inner      *DefaultAIService
}

func NewChatClientChatCompletionService(chatClient chatcompletion.IChatClient, metadata chatcompletion.ChatClientMetadata) *ChatClientChatCompletionService {
	attrs := map[string]interface{}{}
	if metadata.ProviderUri != nil {
		attrs[EndpointKey] = *metadata.ProviderUri
	}
	if metadata.DefaultModelId != nil {
		attrs[ModelIdKey] = *metadata.DefaultModelId
	}
	chat := &ChatClientChatCompletionService{
		chatClient: chatClient,
		inner:      NewDefaultAIService(attrs),
	}
	return chat
}

// GetApiVersion implements IChatCompletionService.
func (c *ChatClientChatCompletionService) GetApiVersion() string {
	return c.inner.GetApiVersion()
}

// GetAttributes implements IChatCompletionService.
func (c *ChatClientChatCompletionService) GetAttributes() map[string]interface{} {
	return c.inner.GetAttributes()
}

// GetEndpoint implements IChatCompletionService.
func (c *ChatClientChatCompletionService) GetEndpoint() string {
	return c.inner.GetEndpoint()
}

// GetModelId implements IChatCompletionService.
func (c *ChatClientChatCompletionService) GetModelId() string {
	return c.inner.GetModelId()
}

// GetChatMessageContents implements IChatCompletionService.
func (c *ChatClientChatCompletionService) GetChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings *PromptExecutionSettings, kernel *Kernel) ([]ChatMessageContent, error) {
	completion, err := c.chatClient.GetResponse(ctx, chatHistory.ToChatMessageList(), c.ToChatOptions(executionSettings, kernel))
	if err != nil {
		return nil, err
	}

	if len(completion.Messages) > 0 {
		for i := 0; i < len(completion.Messages)-1; i++ {
			chatHistory.Add(ToChatMessageContent(completion.Messages[i], completion))
		}

		// Return the last message as the result.
		return []ChatMessageContent{ToChatMessageContent(completion.Messages[len(completion.Messages)-1], completion)}, nil
	}

	return []ChatMessageContent{}, nil
}

// GetStreamingChatMessageContents implements IChatCompletionService.
func (c *ChatClientChatCompletionService) GetStreamingChatMessageContents(ctx context.Context, chatHistory ChatHistory, executionSettings *PromptExecutionSettings, kernel *Kernel) (<-chan StreamingChatMessageContent, <-chan error) {
	fcContents := []aicontents.IAIContent{}
	var role *chatcompletion.ChatRole
	contentCh := make(chan StreamingChatMessageContent)
	errCh := make(chan error, 1)

	go func() {
		defer close(contentCh)
		defer close(errCh)

		responseCh := c.chatClient.GetStreamingResponse(ctx, chatHistory.ToChatMessageList(), c.ToChatOptions(executionSettings, kernel))
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		case update, ok := <-responseCh:
			if !ok {
				return
			}
			if update.Err != nil {
				errCh <- update.Err
				return
			}
			if role == nil {
				role = update.Update.Role
			}

			for _, c := range update.Update.Contents {
				if cc, ok := c.(aicontents.FunctionCallContent); ok {
					fcContents = append(fcContents, cc)
				}
				if cc, ok := c.(aicontents.FunctionResultContent); ok {
					fcContents = append(fcContents, cc)
				}
			}
			contentCh <- c.ToStreamingChatMessageContent(*update.Update)
		}
	}()
	message := &chatcompletion.ChatMessage{
		Role:                 chatcompletion.RoleAssistant,
		Contents:             fcContents,
		AdditionalProperties: map[string]interface{}{},
	}
	if role != nil {
		message.Role = *role
	}
	chatHistory.Add(ToChatMessageContent(*message, nil))
	return contentCh, errCh
}

func (c *ChatClientChatCompletionService) ToStreamingChatMessageContent(update chatcompletion.ChatResponseUpdate) StreamingChatMessageContent {
	content := &StreamingChatMessageContent{
		InnerContent: update.RawRepresentation,
		Metadata:     update.AdditionalProperties,
	}
	if update.ModelId != nil {
		content.ModelId = *update.ModelId
	}
	if update.Role != nil {
		content.Role = CreateAuthorRole(string(*update.Role))
	}

	for _, item := range update.Contents {
		var resultContent StreamingKernelContent
		if tc, ok := item.(aicontents.TextContent); ok {
			resultContent = &StreamingTextContent{Text: tc.Text, ModelId: content.ModelId}
		} else if fcc, ok := item.(aicontents.FunctionCallContent); ok {
			c := &StreamingFunctionCallUpdateContent{CallId: fcc.CallId, Name: fcc.Name, ModelId: content.ModelId}
			if fcc.Arguments != nil {
				data, err := json.Marshal(fcc.Arguments)
				if err == nil {
					c.Arguments = string(data)
				}
			}
			resultContent = c
		}

		if resultContent != nil {
			content.Items.Add(resultContent)
		}
	}

	return *content
}

func (c *ChatClientChatCompletionService) ToChatOptions(settings *PromptExecutionSettings, kernel *Kernel) *chatcompletion.ChatOptions {
	if settings == nil {
		return nil
	}

	options := &chatcompletion.ChatOptions{
		ModelId:              &settings.ModelId,
		AdditionalProperties: map[string]interface{}{},
	}

	extensionCopy := make(map[string]interface{})
	for k, v := range settings.ExtensionData {
		extensionCopy[k] = v
	}

	if v, ok := extensionCopy["temperature"].(float64); ok {
		options.Temperature = &v
		delete(extensionCopy, "temperature")
	}

	if v, ok := extensionCopy["top_p"].(float64); ok {
		options.TopP = &v
		delete(extensionCopy, "top_p")
	}

	if v, ok := extensionCopy["top_k"].(int); ok {
		options.TopK = &v
		delete(extensionCopy, "top_k")
	}

	if v, ok := extensionCopy["seed"].(int64); ok {
		options.Seed = &v
		delete(extensionCopy, "seed")
	}
	if v, ok := extensionCopy["max_tokens"].(int64); ok {
		options.MaxOutputTokens = &v
		delete(extensionCopy, "max_tokens")
	}

	if v, ok := extensionCopy["frequency_penalty"].(float64); ok {
		options.FrequencyPenalty = &v
		delete(extensionCopy, "frequency_penalty")
	}

	if v, ok := extensionCopy["presence_penalty"].(float64); ok {
		options.PresencePenalty = &v
		delete(extensionCopy, "presence_penalty")
	}

	if v, ok := extensionCopy["stop_sequences"].([]string); ok {
		options.StopSequences = v
		delete(extensionCopy, "stop_sequences")
	}

	if v, ok := extensionCopy["response_format"].(string); ok {
		s := chatcompletion.NewChatResponseFormat(v)
		options.ResponseFormat = &s
		delete(extensionCopy, "response_format")
	}
	for k, v := range extensionCopy {
		options.AdditionalProperties[k] = v
	}

	configuration := settings.FunctionChoiceBehavior.GetConfiguration(FunctionChoiceBehaviorConfigurationContext{Kernel: kernel, ChatHistory: ChatHistory{List: *core.NewList[ChatMessageContent]()}})
	if len(configuration.Functions) > 0 {
		a := chatcompletion.NoneMode
		if settings.FunctionChoiceBehavior.GetBehaviorType() == "required" {
			a = chatcompletion.RequireAnyMode
		}
		options.ToolMode = &a
		for _, f := range configuration.Functions {
			options.Tools = append(options.Tools, f)
		}
	}

	return options

}

func TryConvert[T any](value any) (T, bool) {
	var typedValue T
	if value != nil {
		if typedValue, ok := value.(T); ok {
			return typedValue, true
		}

		if jsonValue, ok := value.(map[string]any); ok {
			result, err := json.Marshal(jsonValue)
			if err != nil {
				return typedValue, false
			}
			if err := json.Unmarshal(result, &typedValue); err != nil {
				return typedValue, false
			}
			return typedValue, true
		}

	}

	return typedValue, false
}
