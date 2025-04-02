package types

type CreateMessageRequestParams struct {
	RequestParams
	IncludeContext   *ContextInclusion `json:"includeContext"`
	MaxTokens        *int              `json:"maxTokens"`
	Messages         []SamplingMessage `json:"messages"`
	Metadata         any               `json:"metadata"`
	ModelPreferences ModelPreferences  `json:"modelPreferences"`
	StopSequences    []string          `json:"stopSequences"`
	SystemPrompt     *string           `json:"systemPrompt"`
	Temperature      *float32          `json:"temperature"`
}
