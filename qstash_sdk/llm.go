package qstash

import (
	"context"
)

type LLMService service

func (s *LLMService) CreateChatCompletion(ctx context.Context, request ChatRequest) (*ChatCompletionResponse, error) {
	path := "/llm/v1/chat/completions"
	result := &ChatCompletionResponse{}
	if err := s.client.http.Post(ctx, path, request, result); err != nil {
		return nil, err
	}

	return result, nil
}

// stream,err:= qstashClient.LLM.CreateChatStreamCompletion(ctx, request)
//
//	if err!=nil {
//		doSomething()
//	}
//
// defer stream.Close()
//
// result := make([]ChatCompletionResponse, 0)
//
//	for {
//			if !stream.CanReadStream() {
//				break
//			}
//			response := ChatCompletionResponse{}
//			if err=stream.ReadStream(&response);err!=nil {
//				doSomething()
//			}else{
//				result = append(result, response)
//			}
//		}
func (s *LLMService) CreateChatStreamCompletion(ctx context.Context, request ChatRequest) (*StreamResponse, error) {
	path := "/llm/v1/chat/completions"
	type Alias ChatRequest
	inner := &struct {
		Stream bool `json:"stream"`
		Alias
	}{
		Stream: true,
		Alias:  (Alias)(request),
	}

	return s.client.http.PostStream(ctx, path, inner)
}

type ChatRequest struct {
	Model            string             `json:"model"`
	Messages         []ChatMessage      `json:"messages"`
	FrequencyPenalty float64            `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int     `json:"logit_bias,omitempty"`
	Logprobs         bool               `json:"logprobs,omitempty"`
	TopLogprobs      int                `json:"top_logprobs,omitempty"`
	MaxTokens        int                `json:"max_tokens,omitempty"`
	N                int                `json:"n,omitempty"`
	PresencePenalty  float64            `json:"presence_penalty,omitempty"`
	ResponseFormat   ChatResponseFormat `json:"response_format,omitempty"`
	Seed             float64            `json:"seed,omitempty"`
	Stop             []string           `json:"stop,omitempty"`
	Temperature      float64            `json:"temperature,omitempty"`
	TopP             float64            `json:"top_p,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type ChatResponseFormat struct {
	Type string `json:"type"` //Must be one of text or json_object.
}

type ChatCompletionResponse struct {
	ID                string          `json:"id"`
	Choices           []Choice        `json:"choices"`
	Created           int64           `json:"created"`
	Model             string          `json:"model"`
	SystemFingerprint string          `json:"system_fingerprint,omitempty"`
	Object            string          `json:"object"`
	Usage             UsageStatistics `json:"usage,omitempty"`
}

type Choice struct {
	Delta        *ChatMessage `json:"delta,omitempty"`
	Message      *ChatMessage `json:"message,omitempty"`
	FinishReason string       `json:"finish_reason,omitempty"`
	StopReason   string       `json:"stop_reason,omitempty"`
	Index        int          `json:"index"`
	Logprobs     LogProbs     `json:"logprobs,omitempty"`
}

type LogProbs struct {
	Content []TokenInfo `json:"content,omitempty"`
}

type TokenInfo struct {
	Token       string         `json:"token"`
	Logprob     float64        `json:"logprob"`
	Bytes       []int          `json:"bytes,omitempty"`
	TopLogprobs []TopTokenInfo `json:"top_logprobs,omitempty"`
}

type TopTokenInfo struct {
	Token   string  `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

type UsageStatistics struct {
	CompletionTokens int `json:"completion_tokens,omitempty"`
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}
