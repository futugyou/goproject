package contents

type FunctionResultContent struct {
	MimeType     string         `json:"mimeType"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	CallId       string         `json:"callId"`
	PluginName   string         `json:"pluginName"`
	FunctionName string         `json:"functionName"`
	Result       any            `json:"result"`
	InnerContent any            `json:"-"`
}

func (FunctionResultContent) Type() string {
	return "functionResult"
}

func (f FunctionResultContent) ToChatMessage() ChatMessageContent {
	return ChatMessageContent{
		MimeType: f.MimeType,
		ModelId:  f.ModelId,
		Metadata: f.Metadata,
		Role:     AuthorRoleTool,
		Items: ChatMessageContentItemCollection{
			Items: []KernelContent{f},
		},
	}
}
