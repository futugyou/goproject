package contents

type FunctionCallContent struct {
	MimeType     string         `json:"mimeType"`
	ModelId      string         `json:"modelId"`
	Metadata     map[string]any `json:"metadata"`
	Id           string         `json:"id"`
	PluginName   string         `json:"pluginName"`
	FunctionName string         `json:"functionName"`
	Arguments    map[string]any `json:"arguments"`
	Exception    string         `json:"exception"`
	InnerContent any            `json:"-"`
}

func (FunctionCallContent) Type() string {
	return "functionCall"
}
