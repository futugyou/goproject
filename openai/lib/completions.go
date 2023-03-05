package lib

const completionsPath string = "completions"

type CreateCompletionRequest struct {
	Model            string           `json:"model,omitempty"`
	Prompt           interface{}      `json:"prompt,omitempty"`
	Suffix           string           `json:"suffix,omitempty"`
	MaxTokens        int32            `json:"max_tokens,omitempty"`
	Temperature      float32          `json:"temperature,omitempty"`
	Top_p            float32          `json:"top_p,omitempty"`
	N                int32            `json:"n,omitempty"`
	Stream           bool             `json:"stream"`
	Logprobs         int32            `json:"logprobs,omitempty"`
	Echo             bool             `json:"echo"`
	Stop             interface{}      `json:"stop,omitempty"`
	PresencePenalty  float32          `json:"presence_penalty,omitempty"`
	FrequencyPenalty float32          `json:"frequency_penalty,omitempty"`
	BestOf           int32            `json:"best_of,omitempty"`
	LogitBias        map[string]int32 `json:"logit_bias,omitempty"`
	User             string           `json:"user,omitempty"`
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

func (c *openaiClient) CreateCompletion(request CreateCompletionRequest) *CreateCompletionResponse {
	result := &CreateCompletionResponse{}
	request.Stream = false
	c.httpClient.Post(completionsPath, request, result)
	return result
}

func (c *openaiClient) CreateStreamCompletion(request CreateCompletionRequest) []*CreateCompletionResponse {
	result := make([]*CreateCompletionResponse, 0)
	request.Stream = true

	c.httpClient.PostStream(completionsPath, request)

	defer c.httpClient.Close()

	for {
		if c.httpClient.StreamEnd {
			break
		}

		response := &CreateCompletionResponse{}
		c.httpClient.ReadStream(response)
		if !c.httpClient.StreamEnd {
			result = append(result, response)
		}
	}

	return result
}
