package ollama

import "encoding/json"

type OllamaChatRequest struct {
	Model    string
	Messages []OllamaChatRequestMessage
	Format   json.RawMessage
	Stream   bool
	Tools    []OllamaTool
	Options  *OllamaRequestOptions
}
