package embeddings

type EmbeddingGeneratorMetadata struct {
	ProviderName *string `json:"providerName,omitempty"`
	ProviderUri  *string `json:"providerUri,omitempty"`
	ModelId      *string `json:"modelId,omitempty"`
	Dimensions   *int64  `json:"dimensions,omitempty"`
}
