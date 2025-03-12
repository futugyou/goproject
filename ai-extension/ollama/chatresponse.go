package ollama

type OllamaChatResponse struct {
	Model              *string
	CreatedAt          *string
	TotalDuration      *int64
	LoadDuration       *int64
	DoneReason         *string
	PromptEvalCount    *int64
	PromptEvalDuration *int64
	EvalCount          *int64
	EvalDuration       *int64
	Message            *OllamaChatResponseMessage
	Done               bool
	Error              *string
}
