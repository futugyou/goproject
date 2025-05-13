package embeddings

import "net/url"

type EmbeddingGeneratorMetadata struct {
	ProviderName           *string  `json:"providerName,omitempty"`
	ProviderUri            *url.URL `json:"providerUri,omitempty"`
	DefaultModelId         *string  `json:"defaultModelId,omitempty"`
	DefaultModelDimensions *int64   `json:"defaultModelDimensions,omitempty"`
}
