package protocol

import "encoding/json"

type Tool struct {
	Name        string           `json:"name"`
	Description *string          `json:"description"`
	InputSchema json.RawMessage  `json:"inputSchema"`
	Annotations *ToolAnnotations `json:"annotations"`
}
