package ollama

type OllamaFunctionTool struct {
	Name        string
	Description string
	Parameters  OllamaFunctionToolParameters
}
