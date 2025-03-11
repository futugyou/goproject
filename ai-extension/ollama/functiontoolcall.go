package ollama

type OllamaFunctionToolCall struct {
	Name      string
	Arguments map[string]interface{}
}
