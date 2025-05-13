package embeddings

type EmbeddingGenerationOptions struct {
	ModelId              *string                `json:"modelId,omitempty"`
	Dimensions           *int64                 `json:"dimensions,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}
