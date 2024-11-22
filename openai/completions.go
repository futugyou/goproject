package openai

import "context"

const completionsPath string = "completions"

type CreateCompletionRequest struct {
	Model string `json:"model,omitempty"`
	// The prompt(s) to generate completions for,
	// encoded as a string, array of strings, array of tokens, or array of token arrays.
	Prompt           interface{}      `json:"prompt,omitempty"`
	Suffix           string           `json:"suffix,omitempty"`
	MaxTokens        int32            `json:"max_tokens,omitempty"`
	Temperature      float32          `json:"temperature,omitempty"`
	Top_p            float32          `json:"top_p,omitempty"`
	N                int32            `json:"n,omitempty"`
	Logprobs         int32            `json:"logprobs,omitempty"`
	Echo             bool             `json:"echo"`
	Stop             []string         `json:"stop,omitempty"`
	PresencePenalty  float32          `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32          `json:"frequency_penalty,omitempty"`
	BestOf           int32            `json:"best_of,omitempty"`
	LogitBias        map[string]int32 `json:"logit_bias,omitempty"`
	User             string           `json:"user,omitempty"`
}

type completionRequest struct {
	CreateCompletionRequest
	Stream bool `json:"stream"`
}

type CreateCompletionResponse struct {
	Error   *OpenaiError `json:"error,omitempty"`
	ID      string       `json:"id,omitempty"`
	Object  string       `json:"object,omitempty"`
	Created int32        `json:"created,omitempty"`
	Model   string       `json:"model,omitempty"`
	Choices []Choices    `json:"choices,omitempty"`
	Usage   *Usage       `json:"usage,omitempty"`
}

type CompletionService service

func (c *CompletionService) CreateCompletion(ctx context.Context, request CreateCompletionRequest) *CreateCompletionResponse {
	result := &CreateCompletionResponse{}
	newRequest := completionRequest{
		CreateCompletionRequest: request,
		Stream:                  false,
	}

	c.client.httpClient.Post(ctx, completionsPath, newRequest, result)
	return result
}

// you can read stream in this way.
//
// stream,err:= openai.CreateStreamCompletion(CreateCompletionRequest{})
//
//	if err!=nil {
//		doSomething()
//	}
//
// defer stream.Close()
//
// result := make([]*CreateCompletionResponse, 0)
//
//	for {
//			if !stream.CanReadStream() {
//				break
//			}
//			response := &CreateCompletionResponse{}
//			if err=stream.ReadStream(response);err!=nil {
//				doSomething()
//			}else{
//				result = append(result, response)
//			}
//		}
func (c *CompletionService) CreateStreamCompletion(ctx context.Context, request CreateCompletionRequest) (*StreamResponse, *OpenaiError) {
	newRequest := completionRequest{
		CreateCompletionRequest: request,
		Stream:                  true,
	}

	return c.client.httpClient.PostStream(ctx, completionsPath, newRequest)
}
