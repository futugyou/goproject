package services

import (
	"os"
	"time"

	"github.com/devfeel/mapper"
	openai "github.com/futugyousuzu/go-openai"
)

type EditService struct {
	client *openai.OpenaiClient
}

func NewEditService(client *openai.OpenaiClient) *EditService {
	if client == nil {
		openaikey := os.Getenv("openaikey")
		client = openai.NewClient(openaikey)
	}
	return &EditService{
		client: client,
	}
}

type CreateEditsRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction"`
	N           int32   `json:"n,omitempty"`
	Temperature float32 `json:"temperature,omitempty"`
	Top_p       float32 `json:"top_p,omitempty"`
}

type CreateEditsResponse struct {
	ErrorMessage     string   `json:"error,omitempty"`
	Created          string   `json:"created,omitempty"`
	PromptTokens     int32    `json:"prompt_tokens,omitempty"`
	CompletionTokens int32    `json:"completion_tokens,omitempty"`
	TotalTokens      int32    `json:"total_tokens,omitempty"`
	Texts            []string `json:"texts,omitempty"`
}

func (s *EditService) CreateEdit(request CreateEditsRequest) CreateEditsResponse {
	req := openai.CreateEditsRequest{}
	mapper.AutoMapper(&request, &req)

	response := s.client.Edit.CreateEdits(req)
	result := CreateEditsResponse{}
	if response != nil {
		if response.Error != nil {
			result.ErrorMessage = response.Error.Error()
		}

		if response.Created != 0 {
			result.Created = time.Unix((int64)(response.Created), 0).Format("2006-01-02 15:04:05")
		}

		if response.Usage != nil {
			result.TotalTokens = response.Usage.TotalTokens
			result.CompletionTokens = response.Usage.CompletionTokens
			result.PromptTokens = response.Usage.PromptTokens
		}

		if response.Choices != nil {
			texts := make([]string, 0)
			for i := 0; i < len(response.Choices); i++ {
				texts = append(texts, response.Choices[i].Text)
			}

			result.Texts = texts
		}
	}

	return result
}
