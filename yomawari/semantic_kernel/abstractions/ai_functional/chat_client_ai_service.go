package ai_functional

import (
	"context"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/semantic_kernel/abstractions/services"
)

type ChatClientAIService struct {
	chatClient         chatcompletion.IChatClient
	internalAttributes map[string]interface{}
	defaultAIService   services.DefaultAIService
}

func NewChatClientAIService(chatClient chatcompletion.IChatClient, meta chatcompletion.ChatClientMetadata) *ChatClientAIService {
	s := &ChatClientAIService{
		chatClient:         chatClient,
		internalAttributes: make(map[string]interface{}),
	}
	if meta.DefaultModelId != nil {
		s.internalAttributes["ModelId"] = *meta.DefaultModelId
	}
	if meta.ProviderName != nil {
		s.internalAttributes["ProviderName"] = *meta.ProviderName
	}
	if meta.ProviderUri != nil {
		s.internalAttributes["ProviderUri"] = *meta.ProviderUri
	}

	s.defaultAIService = *services.NewDefaultAIService(s.internalAttributes)
	return s
}

// GetResponse implements chatcompletion.IChatClient.
func (c *ChatClientAIService) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	return c.chatClient.GetResponse(ctx, chatMessages, options)
}

// GetStreamingResponse implements chatcompletion.IChatClient.
func (c *ChatClientAIService) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	return c.chatClient.GetStreamingResponse(ctx, chatMessages, options)
}

// GetApiVersion implements services.IAIService.
func (c *ChatClientAIService) GetApiVersion() string {
	return c.defaultAIService.GetApiVersion()
}

// GetAttributes implements services.IAIService.
func (c *ChatClientAIService) GetAttributes() map[string]interface{} {
	return c.defaultAIService.GetAttributes()
}

// GetEndpoint implements services.IAIService.
func (c *ChatClientAIService) GetEndpoint() string {
	return c.defaultAIService.GetEndpoint()
}

// GetModelId implements services.IAIService.
func (c *ChatClientAIService) GetModelId() string {
	return c.defaultAIService.GetModelId()
}
