package chatcompletion

import "github.com/futugyou/ai-extension/abstractions"

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
	Tools                []abstractions.AITool  `json:"-"`
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

func (o *ChatOptions) Clone() *ChatOptions {
	if o == nil {
		return &ChatOptions{}
	}

	options := ChatOptions{
		ChatThreadId:         o.ChatThreadId,
		Temperature:          o.Temperature,
		MaxOutputTokens:      o.MaxOutputTokens,
		TopP:                 o.TopP,
		TopK:                 o.TopK,
		FrequencyPenalty:     o.FrequencyPenalty,
		PresencePenalty:      o.PresencePenalty,
		Seed:                 o.Seed,
		ResponseFormat:       o.ResponseFormat,
		ModelId:              o.ModelId,
		ToolMode:             o.ToolMode,
		AdditionalProperties: o.AdditionalProperties,
		StopSequences:        o.StopSequences,
		Tools:                o.Tools,
	}

	additionalProperties := map[string]interface{}{}
	for k, v := range o.AdditionalProperties {
		additionalProperties[k] = v
	}

	options.AdditionalProperties = additionalProperties

	return &options
}
