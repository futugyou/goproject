package chatcompletion

type ChatClientMetadata struct {
	ProviderName *string `json:"providerName,omitempty"`
	ProviderUri  *string `json:"providerUri,omitempty"`
	ModelId      *string `json:"modelId,omitempty"`
}
