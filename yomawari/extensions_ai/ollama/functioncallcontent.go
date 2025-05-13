package ollama

import "encoding/json"

type OllamaFunctionCallContent struct {
	CallId    *string
	Name      *string
	Arguments json.RawMessage
}
