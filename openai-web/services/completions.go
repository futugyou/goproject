package services

import (
	"time"

	openai "github.com/futugyousuzu/go-openai"

	"github.com/beego/beego/v2/core/config"
	"github.com/devfeel/mapper"
)

type CompletionService struct {
}

type CompletionModel struct {
	Model            string   `json:"model,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	MaxTokens        int32    `json:"max_tokens,omitempty"`
	Top_p            float32  `json:"top_p,omitempty"`
	FrequencyPenalty float32  `json:"frequency_penalty,omitempty"`
	PresencePenalty  float32  `json:"presence_penalty,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	BestOf           int32    `json:"best_of,omitempty"`
}

type CreateCompletionRequest struct {
	CompletionModel
}

type CreateCompletionResponse struct {
	ErrorMessage     string   `json:"error,omitempty"`
	Created          string   `json:"created,omitempty"`
	PromptTokens     int32    `json:"prompt_tokens,omitempty"`
	CompletionTokens int32    `json:"completion_tokens,omitempty"`
	TotalTokens      int32    `json:"total_tokens,omitempty"`
	Texts            []string `json:"texts,omitempty"`
}

func (s *CompletionService) CreateCompletion(request CreateCompletionRequest) CreateCompletionResponse {
	openaikey, _ := config.String("openaikey")
	client := openai.NewClient(openaikey)
	req := openai.CreateCompletionRequest{}
	mapper.AutoMapper(&request.CompletionModel, &req)
	apiresult := client.CreateCompletion(req)

	result := CreateCompletionResponse{}
	if apiresult != nil {
		if apiresult.Error != nil {
			result.ErrorMessage = apiresult.Error.Error()
		}

		if apiresult.Created != 0 {
			result.Created = time.Unix((int64)(apiresult.Created), 0).Format(time.DateTime)
		}

		if apiresult.Usage != nil {
			result.TotalTokens = apiresult.Usage.TotalTokens
			result.CompletionTokens = apiresult.Usage.CompletionTokens
			result.PromptTokens = apiresult.Usage.PromptTokens
		}

		if apiresult.Choices != nil {
			texts := make([]string, 0)
			for i := 0; i < len(apiresult.Choices); i++ {
				texts = append(texts, apiresult.Choices[i].Text)
			}

			result.Texts = texts
		}
	}

	return result
}
