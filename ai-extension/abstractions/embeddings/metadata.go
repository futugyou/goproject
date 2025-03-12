package embeddings

import "net/url"

type EmbeddingGeneratorMetadata struct {
	ProviderName *string  `json:"providerName,omitempty"`
	ProviderUri  *url.URL `json:"providerUri,omitempty"`
	ModelId      *string  `json:"modelId,omitempty"`
	Dimensions   *int64   `json:"dimensions,omitempty"`
}
