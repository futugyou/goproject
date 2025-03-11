package ollama

type OllamaChatRequestMessage struct {
	Role    string
	Content *string
	Images  []string
}
