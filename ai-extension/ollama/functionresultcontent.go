package ollama

import "encoding/json"

type OllamaFunctionResultContent struct {
	CallId *string
	Result json.RawMessage
}
