package ollama

type OllamaChatResponse struct {
	Model              *string
	CreatedAt          *string
	TotalDuration      *int64
	LoadDuration       *int64
	DoneReason         *string
	PromptEvalCount    *int
	PromptEvalDuration *int64
	EvalCount          *int
	EvalDuration       *int64
	Message            *OllamaChatResponseMessage
	Done               bool
	Error              *string
}
