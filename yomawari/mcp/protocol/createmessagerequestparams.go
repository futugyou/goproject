package protocol

type CreateMessageRequestParams struct {
	RequestParams    `json:",inline"`
	IncludeContext   *ContextInclusion `json:"includeContext"`
	MaxTokens        *int64            `json:"maxTokens"`
	Messages         []SamplingMessage `json:"messages"`
	Metadata         any               `json:"metadata"`
	ModelPreferences ModelPreferences  `json:"modelPreferences"`
	StopSequences    []string          `json:"stopSequences"`
	SystemPrompt     *string           `json:"systemPrompt"`
	Temperature      *float64          `json:"temperature"`
}
