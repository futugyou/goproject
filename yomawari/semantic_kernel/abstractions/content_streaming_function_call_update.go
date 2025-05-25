package abstractions

import "encoding/base64"

type StreamingFunctionCallUpdateContent struct {
	ChoiceIndex       int            `json:"choiceIndex"`
	ModelId           string         `json:"modelId"`
	Metadata          map[string]any `json:"metadata"`
	InnerContent      any            `json:"-"`
	CallId            string         `json:"callId"`
	Name              string         `json:"name"`
	Arguments         string         `json:"arguments"`
	FunctionCallIndex int            `json:"functionCallIndex"`
	RequestIndex      int            `json:"requestIndex"`
}

func (StreamingFunctionCallUpdateContent) Type() string {
	return "streaming-function-call-update"
}

func (c StreamingFunctionCallUpdateContent) ToString() string {
	return "StreamingFunctionCallUpdateContent"
}

func (c StreamingFunctionCallUpdateContent) ToByteArray() []byte {
	r, _ := base64.URLEncoding.DecodeString(c.ToString())
	return r
}

func (c StreamingFunctionCallUpdateContent) Hash() string {
	return c.ToString()
}
