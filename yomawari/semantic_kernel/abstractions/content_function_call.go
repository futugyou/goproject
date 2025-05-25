package abstractions

import (
	"context"
)

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

func (f FunctionCallContent) ToString() string {
	return f.FunctionName
}

func (c FunctionCallContent) Hash() string {
	return c.ToString()
}

func (f FunctionCallContent) GetInnerContent() any {
	return f.InnerContent
}

func (FunctionCallContent) InvokeAsync(ctx context.Context, kernel Kernel) (*FunctionResultContent, error) {
	// TODO
	return nil, nil
}
