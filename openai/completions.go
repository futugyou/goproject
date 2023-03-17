package openai

import (
	"golang.org/x/exp/slices"

	e "github.com/futugyousuzu/go-openai/internal"
)

const completionsPath string = "completions"

// the model in https://platform.openai.com/docs/models/model-endpoint-compatibility
// is not same as https://platform.openai.com/playground.
// i think i need add those two code models.
var supportedCompletionModel = []string{
	Text_davinci_003,
	Text_davinci_002,
	Text_curie_001,
	Text_babbage_001,
	Text_ada_001,
	GPT3_davinci,
	GPT3_curie,
	GPT3_babbage,
	GPT3_ada,
	Code_davinci_002,
	Code_cushman_001,
}

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
	Stream           bool             `json:"stream"`
	Logprobs         int32            `json:"logprobs,omitempty"`
	Echo             bool             `json:"echo"`
	Stop             []string         `json:"stop,omitempty"`
	PresencePenalty  float32          `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32          `json:"frequency_penalty,omitempty"`
	BestOf           int32            `json:"best_of,omitempty"`
	LogitBias        map[string]int32 `json:"logit_bias,omitempty"`
	User             string           `json:"user,omitempty"`
}

type CreateCompletionResponse struct {
	Error   *e.OpenaiError `json:"error,omitempty"`
	ID      string         `json:"id,omitempty"`
	Object  string         `json:"object,omitempty"`
	Created int32          `json:"created,omitempty"`
	Model   string         `json:"model,omitempty"`
	Choices []Choices      `json:"choices,omitempty"`
	Usage   *Usage         `json:"usage,omitempty"`
}

func (c *openaiClient) CreateCompletion(request CreateCompletionRequest) *CreateCompletionResponse {
	result := &CreateCompletionResponse{}
	request.Stream = false
	err := validateCompletionModel(request.Model)
	if err != nil {
		result.Error = err
		return result
	}

	c.httpClient.Post(completionsPath, request, result)
	return result
}

func (c *openaiClient) CreateStreamCompletion(request CreateCompletionRequest) []*CreateCompletionResponse {
	result := make([]*CreateCompletionResponse, 0)
	request.Stream = true
	err := validateCompletionModel(request.Model)
	if err != nil {
		result = append(result, &CreateCompletionResponse{Error: err})
		return result
	}

	c.httpClient.PostStream(completionsPath, request)

	defer c.httpClient.Close()

	for {
		if !c.httpClient.CanReadStream() {
			break
		}

		response := &CreateCompletionResponse{}
		c.httpClient.ReadStream(response)
		if c.httpClient.CanReadStream() {
			result = append(result, response)
		}
	}

	return result
}

func validateCompletionModel(model string) *e.OpenaiError {
	if len(model) == 0 || !slices.Contains(supportedCompletionModel, model) {
		return e.UnsupportedTypeError("Model", model, supportedCompletionModel)
	}

	return nil
}
