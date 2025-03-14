package ollama

type OllamaTool struct {
	Type     string
	Function OllamaFunctionTool
}
