package contents

import "encoding/json"

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

func (f FunctionResultContent) ToString() string {
	d, err := json.Marshal(f.Result)
	if err != nil {
		return ""
	}
	return string(d)
}

func (f FunctionResultContent) GetInnerContent() any {
	return f.InnerContent
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
