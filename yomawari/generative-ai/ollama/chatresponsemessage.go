package ollama

type OllamaChatResponseMessage struct {
	Role      string
	Content   string
	ToolCalls []OllamaToolCall
}
