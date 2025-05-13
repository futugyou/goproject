package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/google/uuid"
)

var schemalessJsonResponseFormatValue = json.RawMessage([]byte(`"json"`))
var headerData []byte = []byte("data: ")

const endTag string = "[DONE]"

type OllamaChatClient struct {
	metadata        *chatcompletion.ChatClientMetadata
	apiChatEndpoint url.URL
	httpClient      *http.Client
}

func NewOllamaChatClient(modelId *string, endpoint url.URL, httpClient *http.Client) *OllamaChatClient {
	var _ chatcompletion.IChatClient = (*OllamaChatClient)(nil)

	name := "ollama"
	apiChatEndpoint := endpoint.JoinPath("api/chat")
	if httpClient == nil {
		httpClient = SharedClient
	}

	return &OllamaChatClient{
		metadata: &chatcompletion.ChatClientMetadata{
			ProviderName:   &name,
			ProviderUri:    &endpoint,
			DefaultModelId: modelId,
		},
		apiChatEndpoint: *apiChatEndpoint,
		httpClient:      httpClient,
	}
}

func (client *OllamaChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	path := client.apiChatEndpoint.String()
	request := ToOllamaChatRequest(chatMessages, options, false)
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("POST", path, body)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var chatResponse OllamaChatResponse
	err = json.NewDecoder(resp.Body).Decode(&chatResponse)
	if err != nil {
		return nil, err
	}

	result := ToChatResponse(chatResponse)
	return result, nil
}

func (client *OllamaChatClient) GetStreamingResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	// Create the channel for streaming responses
	response := make(chan chatcompletion.ChatStreamingResponse)

	// Prepare the request
	path := client.apiChatEndpoint.String()
	request := ToOllamaChatRequest(chatMessages, options, true)
	payloadBytes, _ := json.Marshal(request)
	body := bytes.NewReader(payloadBytes)

	req, _ := http.NewRequest("POST", path, body)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Connection", "keep-alive")

	// Perform the HTTP request
	resp, err := client.httpClient.Do(req)
	if err != nil {
		// Handle error and close the channel
		go func() {
			response <- chatcompletion.ChatStreamingResponse{
				Update: nil,
				Err:    err,
			}
			close(response)
		}()
		return response
	}
	defer resp.Body.Close()

	// Initialize the reader for the response body
	reader := bufio.NewReader(resp.Body)

	// Handle streaming data
	go func() {
		defer close(response) // Ensure channel is closed at the end

		for {
			select {
			case <-ctx.Done():
				// Handle context cancellation
				response <- chatcompletion.ChatStreamingResponse{
					Update: nil,
					Err:    ctx.Err(),
				}
				return
			default:
				// Read each line of the stream
				line, err := reader.ReadBytes('\n')
				if err != nil {
					// Handle read errors
					response <- chatcompletion.ChatStreamingResponse{
						Update: nil,
						Err:    err,
					}
					return
				}

				line = bytes.TrimSpace(line)
				if bytes.HasPrefix(line, headerData) {
					line = bytes.TrimPrefix(line, headerData)
					responseStr := string(line)

					// Check for end of response
					if responseStr == endTag {
						return
					}

					// Parse the response JSON
					chatResponse := OllamaChatResponse{}
					if err := json.Unmarshal(line, &chatResponse); err != nil {
						response <- chatcompletion.ChatStreamingResponse{
							Update: nil,
							Err:    err,
						}
						return
					}

					// Construct the update response
					reason := ToFinishReason(chatResponse)
					responseId := uuid.New().String()
					update := chatcompletion.ChatResponseUpdate{
						CreatedAt:    StringToTimePtr(chatResponse.CreatedAt, time.RFC3339),
						ResponseId:   &responseId,
						FinishReason: &reason,
						ModelId:      chatResponse.Model,
						Contents:     []contents.IAIContent{},
					}

					// Process the message content and tool calls
					if chatResponse.Message != nil {
						role := chatcompletion.StringToChatRole(chatResponse.Message.Role)
						update.Role = &role
						for _, tool := range chatResponse.Message.ToolCalls {
							if tool.Function != nil {
								update.Contents = append(update.Contents, ToFunctionCallContent(*tool.Function))
							}
						}

						if len(chatResponse.Message.Content) > 0 || len(update.Contents) == 0 {
							update.Contents = append(update.Contents, contents.TextContent{
								Text: chatResponse.Message.Content,
							})
						}
					}

					// Add usage content if available
					useage := ParseOllamaChatResponseUsage(chatResponse)
					if useage != nil {
						update.Contents = append(update.Contents, contents.NewUsageContent(*useage))
					}

					// Send the update
					response <- chatcompletion.ChatStreamingResponse{
						Update: &update,
						Err:    nil,
					}
				}
			}
		}
	}()

	// Return the response channel
	return response
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

		request.Tools = ToOllamaTools(options.Tools, options.ToolMode)

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

func ToOllamaTools(aITool []abstractions.AITool, chatToolMode *chatcompletion.ChatToolMode) []OllamaTool {
	if chatToolMode != nil && *chatToolMode == chatcompletion.NoneMode {
		return []OllamaTool{}
	}

	result := make([]OllamaTool, len(aITool))
	for i := 0; i < len(result); i++ {
		result[i] = OllamaTool{
			Type: "function",
			// TODO: need json schema for AITool
			Function: OllamaFunctionTool{
				Name:        aITool[i].GetName(),
				Description: aITool[i].GetDescription(),
			},
		}
	}
	return result
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

func FromOllamaMessage(message OllamaChatResponseMessage) chatcompletion.ChatMessage {
	conts := []contents.IAIContent{}
	for _, tool := range message.ToolCalls {
		if tool.Function != nil {
			conts = append(conts, ToFunctionCallContent(*tool.Function))
		}
	}

	if len(message.Content) > 0 || len(conts) == 0 {
		conts = append(conts, contents.TextContent{
			Text: message.Content,
		})
	}

	return chatcompletion.ChatMessage{
		Role:     chatcompletion.StringToChatRole(message.Role),
		Contents: conts,
	}
}

func ToFunctionCallContent(ollamaFunctionToolCall OllamaFunctionToolCall) contents.IAIContent {
	return contents.FunctionCallContent{
		CallId:    uuid.New().String(),
		Name:      ollamaFunctionToolCall.Name,
		Arguments: ollamaFunctionToolCall.Arguments,
	}
}

func ParseOllamaChatResponseUsage(chatResponse OllamaChatResponse) *abstractions.UsageDetails {
	metadata := map[string]int64{}
	TransferNanosecondsTime(chatResponse, func(response OllamaChatResponse) *int64 { return response.LoadDuration }, "load_duration", &metadata)
	TransferNanosecondsTime(chatResponse, func(response OllamaChatResponse) *int64 { return response.TotalDuration }, "total_duration", &metadata)
	TransferNanosecondsTime(chatResponse, func(response OllamaChatResponse) *int64 { return response.PromptEvalDuration }, "prompt_eval_duration", &metadata)
	TransferNanosecondsTime(chatResponse, func(response OllamaChatResponse) *int64 { return response.EvalDuration }, "eval_duration", &metadata)

	if len(metadata) > 0 || chatResponse.PromptEvalCount != nil || chatResponse.EvalCount != nil {
		var p int64 = 0
		if chatResponse.PromptEvalCount != nil {
			p = *chatResponse.PromptEvalCount
		}

		if chatResponse.EvalCount != nil {
			p = p + *chatResponse.EvalCount + p
		}

		return &abstractions.UsageDetails{
			InputTokenCount:      chatResponse.PromptEvalCount,
			OutputTokenCount:     chatResponse.EvalCount,
			TotalTokenCount:      &p,
			AdditionalProperties: metadata,
		}
	}

	return nil
}

func ToFinishReason(chatResponse OllamaChatResponse) chatcompletion.ChatFinishReason {
	if chatResponse.DoneReason == nil {
		return chatcompletion.ReasonUnknown
	}
	switch *chatResponse.DoneReason {
	case "length":
		return chatcompletion.ReasonLength
	case "stop":
		return chatcompletion.ReasonStop
	default:
		return chatcompletion.ReasonUnknown
	}
}

func ToChatResponse(chatResponse OllamaChatResponse) *chatcompletion.ChatResponse {
	message := chatcompletion.ChatMessage{}
	if chatResponse.Message != nil {
		message = FromOllamaMessage(*chatResponse.Message)
	}

	reason := ToFinishReason(chatResponse)
	result := chatcompletion.NewChatResponse(nil, &message)
	result.ResponseId = chatResponse.CreatedAt
	result.ModelId = chatResponse.Model
	result.CreatedAt = StringToTimePtr(chatResponse.CreatedAt, time.RFC3339)
	result.FinishReason = &reason
	result.Usage = ParseOllamaChatResponseUsage(chatResponse)

	return result
}
