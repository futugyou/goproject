package chatcompletion

type ChatOptions struct {
	ChatThreadId         *string                `json:"chatThreadId,omitempty"`
	Temperature          *float32               `json:"temperature,omitempty"`
	MaxOutputTokens      *int                   `json:"maxOutputTokens,omitempty"`
	TopP                 *float32               `json:"topP,omitempty"`
	TopK                 *int                   `json:"topK,omitempty"`
	FrequencyPenalty     *float32               `json:"frequencyPenalty,omitempty"`
	PresencePenalty      *float32               `json:"presencePenalty,omitempty"`
	Seed                 *int64                 `json:"seed,omitempty"`
	ResponseFormat       *ChatResponseFormat    `json:"responseFormat,omitempty"`
	ModelId              *string                `json:"modelId,omitempty"`
	StopSequences        []string               `json:"stopSequences,omitempty"`
	ToolMode             *ChatToolMode          `json:"toolMode,omitempty"`
	Tools                []AITool               `json:"tools,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additionalProperties,omitempty"`
}

type ChatResponseFormat string

const (
	TextFormat ChatResponseFormat = "text"
	JsonFormat ChatResponseFormat = "json"
)

type ChatToolMode string

const (
	AutoMode       ChatToolMode = "auto"
	ManualMode     ChatToolMode = "manual"
	RequireAnyMode ChatToolMode = "requireAny"
	NoneMode       ChatToolMode = "none"
)

type AITool struct {
	// Define the structure based on the specific tool's fields
}
