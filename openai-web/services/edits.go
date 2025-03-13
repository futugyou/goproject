package services

import (
	"context"

	openai "github.com/openai/openai-go"
)

// EditService is deprecated.
// Deprecated: This service is no longer to use
type EditService struct {
	client *openai.Client
}

func NewEditService(client *openai.Client) *EditService {
	return &EditService{
		client: client,
	}
}

type CreateEditsRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction"`
	N           int64   `json:"n,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	Top_p       float64 `json:"top_p,omitempty"`
}

type CreateEditsResponse struct {
	ErrorMessage     string   `json:"error,omitempty"`
	Created          string   `json:"created,omitempty"`
	PromptTokens     int64    `json:"prompt_tokens,omitempty"`
	CompletionTokens int64    `json:"completion_tokens,omitempty"`
	TotalTokens      int64    `json:"total_tokens,omitempty"`
	Texts            []string `json:"texts,omitempty"`
}

func (s *EditService) CreateEdit(ctx context.Context, request CreateEditsRequest) CreateEditsResponse {
	result := CreateEditsResponse{}
	return result
}
