package ollama

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
	"github.com/futugyou/ai-extension/abstractions/contents"
)

var schemalessJsonResponseFormatValue = json.RawMessage([]byte(`"json"`))

type OllamaChatClient struct {
	metadata        *chatcompletion.ChatClientMetadata
	apiChatEndpoint *url.URL
	httpClient      *http.Client
}

func NewOllamaChatClient(modelId *string, endpoint url.URL, httpClient *http.Client) *OllamaChatClient {
	name := "ollama"
	apiChatEndpoint := endpoint.JoinPath("api/chat")
	if httpClient == nil {
		httpClient = SharedClient
	}

	return &OllamaChatClient{
		metadata: &chatcompletion.ChatClientMetadata{
			ProviderName: &name,
			ProviderUri:  &endpoint,
			ModelId:      modelId,
		},
		apiChatEndpoint: apiChatEndpoint,
		httpClient:      httpClient,
	}
}

func (client *OllamaChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	// TODO: Implement
	return nil, nil
}

func (client *OllamaChatClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	// TODO: Implement
	return nil
}

func ToOllamaChatRequest(messages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions, stream bool) *OllamaChatRequest {
	request := &OllamaChatRequest{
		Options:  &OllamaRequestOptions{},
		Stream:   stream,
		Model:    "",
		Messages: ToOllamaChatRequestMessages(messages),
	}

	if options != nil {
		request.Format = ToOllamaChatResponseFormat(options.ResponseFormat)
		if options.ModelId != nil {
			request.Model = *options.ModelId
		}

		transferMetadataValue(options, request, "logits_all", func(options *OllamaRequestOptions, value bool) {
			options.LogitsAll = &value
		})
		transferMetadataValue(options, request, "low_vram", func(options *OllamaRequestOptions, value bool) {
			options.LowVRAM = &value
		})
		transferMetadataValue(options, request, "main_gpu", func(options *OllamaRequestOptions, value int) {
			options.MainGPU = &value
		})
		transferMetadataValue(options, request, "min_p", func(options *OllamaRequestOptions, value float64) {
			options.MinP = &value
		})
		transferMetadataValue(options, request, "mirostat", func(options *OllamaRequestOptions, value int) {
			options.Mirostat = &value
		})
		transferMetadataValue(options, request, "mirostat_eta", func(options *OllamaRequestOptions, value float64) {
			options.MirostatEta = &value
		})
		transferMetadataValue(options, request, "mirostat_tau", func(options *OllamaRequestOptions, value float64) {
			options.MirostatTau = &value
		})
		transferMetadataValue(options, request, "num_batch", func(options *OllamaRequestOptions, value int) {
			options.NumBatch = &value
		})
		transferMetadataValue(options, request, "num_ctx", func(options *OllamaRequestOptions, value int) {
			options.NumCtx = &value
		})
		transferMetadataValue(options, request, "num_gpu", func(options *OllamaRequestOptions, value int) {
			options.NumGPU = &value
		})
		transferMetadataValue(options, request, "num_keep", func(options *OllamaRequestOptions, value int) {
			options.NumKeep = &value
		})
		transferMetadataValue(options, request, "num_thread", func(options *OllamaRequestOptions, value int) {
			options.NumThread = &value
		})
		transferMetadataValue(options, request, "numa", func(options *OllamaRequestOptions, value bool) {
			options.NUMA = &value
		})
		transferMetadataValue(options, request, "penalize_newline", func(options *OllamaRequestOptions, value bool) {
			options.PenalizeNewline = &value
		})
		transferMetadataValue(options, request, "repeat_last_n", func(options *OllamaRequestOptions, value int) {
			options.RepeatLastN = &value
		})
		transferMetadataValue(options, request, "repeat_penalty", func(options *OllamaRequestOptions, value float64) {
			options.RepeatPenalty = &value
		})
		transferMetadataValue(options, request, "tfs_z", func(options *OllamaRequestOptions, value float64) {
			options.TFSZ = &value
		})
		transferMetadataValue(options, request, "typical_p", func(options *OllamaRequestOptions, value float64) {
			options.TypicalP = &value
		})
		transferMetadataValue(options, request, "use_mlock", func(options *OllamaRequestOptions, value bool) {
			options.UseMLock = &value
		})
		transferMetadataValue(options, request, "use_mmap", func(options *OllamaRequestOptions, value bool) {
			options.UseMMap = &value
		})
		transferMetadataValue(options, request, "vocab_only", func(options *OllamaRequestOptions, value bool) {
			options.VocabOnly = &value
		})

		if options.FrequencyPenalty != nil {
			request.Options.FrequencyPenalty = options.FrequencyPenalty
		}

		if options.MaxOutputTokens != nil {
			request.Options.NumPredict = options.MaxOutputTokens
		}

		if options.PresencePenalty != nil {
			request.Options.PresencePenalty = options.PresencePenalty
		}

		if options.StopSequences != nil {
			request.Options.Stop = options.StopSequences
		}

		if options.Temperature != nil {
			request.Options.Temperature = options.Temperature
		}

		if options.TopP != nil {
			request.Options.TopP = options.TopP
		}

		if options.TopK != nil {
			request.Options.TopK = options.TopK
		}

		if options.Seed != nil {
			request.Options.Seed = options.Seed
		}
	}

	return request
}

func ToOllamaChatRequestMessages(messages []chatcompletion.ChatMessage) []OllamaChatRequestMessage {
	response := []OllamaChatRequestMessage{}
	var currentTextMessage *OllamaChatRequestMessage = nil

	for _, message := range messages {
		for _, content := range message.Contents {
			if con, ok := content.(*contents.DataContent); ok && con.MediaTypeStartsWith("image") {
				if currentTextMessage != nil {
					currentTextMessage.Images = append(currentTextMessage.Images, base64.StdEncoding.EncodeToString(con.Data))
				} else {
					response = append(response, OllamaChatRequestMessage{
						Images: []string{base64.StdEncoding.EncodeToString(con.Data)},
						Role:   string(message.Role),
					})
				}
			} else {
				if currentTextMessage != nil {
					response = append(response, *currentTextMessage)
					currentTextMessage = nil
				}

				switch cont := content.(type) {
				case *contents.TextContent:
					text := cont.Text
					currentTextMessage = &OllamaChatRequestMessage{
						Content: &text,
						Role:    string(message.Role),
					}

				case *contents.FunctionCallContent:
					arguments, _ := json.Marshal(cont.Arguments)
					callContent := OllamaFunctionCallContent{
						CallId:    &cont.CallId,
						Name:      &cont.Name,
						Arguments: arguments,
					}
					jsonString, _ := json.Marshal(callContent)

					response = append(response, OllamaChatRequestMessage{
						Content: Ptr(string(jsonString)),
						Role:    "assistant",
					})

				case *contents.FunctionResultContent:
					result, _ := json.Marshal(cont.Result)
					callContent := OllamaFunctionResultContent{
						CallId: &cont.CallId,
						Result: result,
					}
					jsonString, _ := json.Marshal(callContent)

					response = append(response, OllamaChatRequestMessage{
						Content: Ptr(string(jsonString)),
						Role:    "tool",
					})
				}
			}
		}
	}

	if currentTextMessage != nil {
		response = append(response, *currentTextMessage)
	}

	return response
}

func Ptr[T any](v T) *T {
	return &v
}

func ToOllamaChatResponseFormat(chatResponseFormat *chatcompletion.ChatResponseFormat) json.RawMessage {
	if chatResponseFormat == nil {
		return schemalessJsonResponseFormatValue
	}

	format := strings.ToLower(string(*chatResponseFormat))
	switch format {
	case "text":
		return json.RawMessage([]byte(`"text"`))
	case "json":
		return json.RawMessage([]byte(`"json"`))
	}

	return schemalessJsonResponseFormatValue
}

func transferMetadataValue[T any](
	options *chatcompletion.ChatOptions,
	request *OllamaChatRequest,
	propertyName string,
	setOption func(*OllamaRequestOptions, T),
) {
	if options == nil || request == nil || request.Options == nil {
		return
	}

	if value, ok := options.AdditionalProperties[propertyName]; ok {
		if condition, ok := value.(T); ok {
			setOption(request.Options, condition)
		}
	}
}
