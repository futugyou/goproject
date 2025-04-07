package chatcompletion

import "net/url"

type ChatClientMetadata struct {
	ProviderName *string  `json:"providerName,omitempty"`
	ProviderUri  *url.URL `json:"providerUri,omitempty"`
	ModelId      *string  `json:"modelId,omitempty"`
}
