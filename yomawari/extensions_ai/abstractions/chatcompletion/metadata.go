package chatcompletion

import "net/url"

type ChatClientMetadata struct {
	ProviderName   *string  `json:"providerName,omitempty"`
	ProviderUri    *url.URL `json:"providerUri,omitempty"`
	DefaultModelId *string  `json:"modelId,omitempty"`
}

func NewChatClientMetadata(providerName *string, providerUri *url.URL, defaultModelId *string) *ChatClientMetadata {
	return &ChatClientMetadata{
		ProviderName:   providerName,
		ProviderUri:    providerUri,
		DefaultModelId: defaultModelId,
	}
}
