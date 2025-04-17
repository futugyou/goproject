package shared

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"sync"

	"github.com/futugyou/yomawari/extensions-ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions-ai/abstractions/contents"
	"github.com/futugyou/yomawari/mcp"
	"github.com/futugyou/yomawari/mcp/protocol/messages"
	"github.com/futugyou/yomawari/mcp/protocol/transport"
	"github.com/futugyou/yomawari/mcp/protocol/types"
)

var _ IMcpEndpoint = (*BaseMcpEndpoint)(nil)

type BaseMcpEndpoint struct {
	mu            sync.Mutex
	disposed      bool
	session       *McpSession
	sessionCts    context.CancelFunc
	messageTask   <-chan struct{}
	reqHandlers   *RequestHandlers
	notifHandlers *NotificationHandlers
	endpointName  string
}

func (e *BaseMcpEndpoint) GetMcpSession() *McpSession {
	return e.session
}

// NotifyProgress implements IMcpEndpoint.
func (e *BaseMcpEndpoint) NotifyProgress(ctx context.Context, progressToken messages.ProgressToken, progress messages.ProgressNotificationValue) error {
	p := messages.ProgressNotification{ProgressToken: &progressToken, Progress: &progress}
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}
	notification := messages.NewJsonRpcNotification(messages.NotificationMethods_ProgressNotification, data)
	return e.SendNotification(ctx, *notification)
}

// SendNotification implements IMcpEndpoint.
func (e *BaseMcpEndpoint) SendNotification(ctx context.Context, notification messages.JsonRpcNotification) error {
	return e.SendMessage(ctx, &notification)
}

func NewBaseMcpEndpoint() *BaseMcpEndpoint {
	return &BaseMcpEndpoint{
		reqHandlers:   NewRequestHandlers(),
		notifHandlers: NewNotificationHandlers(),
		endpointName:  "",
	}
}

func (e *BaseMcpEndpoint) GetEndpointName() string {
	return e.endpointName
}

func (e *BaseMcpEndpoint) GetMessageProcessingTask() <-chan struct{} {
	return e.messageTask
}

func (e *BaseMcpEndpoint) GetRequestHandlers() *RequestHandlers {
	return e.reqHandlers
}

func (e *BaseMcpEndpoint) GetNotificationHandlers() *NotificationHandlers {
	return e.notifHandlers
}

func (e *BaseMcpEndpoint) InitializeSession(transport transport.ITransport, isServer bool) {
	e.session = NewMcpSession(isServer, transport, e.endpointName, e.reqHandlers, e.notifHandlers)
}

func (e *BaseMcpEndpoint) StartSession(ctx context.Context, transport transport.ITransport) {
	childCtx, cancel := context.WithCancel(ctx)
	e.sessionCts = cancel

	done := make(chan struct{})
	e.messageTask = done

	go func() {
		defer close(done)
		e.session.ProcessMessages(childCtx)
	}()
}

func (e *BaseMcpEndpoint) CancelSession() {
	if e != nil && e.sessionCts != nil {
		e.sessionCts()
	}
}

func (e *BaseMcpEndpoint) SendRequest(ctx context.Context, req *messages.JsonRpcRequest) (*messages.JsonRpcResponse, error) {
	if e == nil || e.session == nil {
		return nil, errors.New("session not initialized")
	}
	return e.session.SendRequest(ctx, req)
}

func (e *BaseMcpEndpoint) SendMessage(ctx context.Context, msg messages.IJsonRpcMessage) error {
	if e == nil || e.session == nil {
		return errors.New("session not initialized")
	}
	return e.session.SendMessage(ctx, msg)
}

func (e *BaseMcpEndpoint) RegisterNotificationHandler(method string, handler NotificationHandler) *RegistrationHandle {
	if e.session == nil {
		return nil
	}
	return e.session.RegisterNotificationHandler(method, handler)
}

func (e *BaseMcpEndpoint) Dispose(ctx context.Context) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.disposed {
		return nil
	}
	e.disposed = true
	return e.disposeUnsynchronized(ctx)
}

func (e *BaseMcpEndpoint) disposeUnsynchronized(ctx context.Context) error {
	if e.sessionCts != nil {
		e.sessionCts()
	}

	if e.messageTask != nil {
		select {
		case <-e.messageTask:
		case <-ctx.Done():
		}
	}

	e.session.Dispose()
	return nil
}

func (e *BaseMcpEndpoint) RequestSampling(ctx context.Context, request types.CreateMessageRequestParams) (*types.CreateMessageResult, error) {
	req := messages.NewJsonRpcRequest(messages.RequestMethods_SamplingCreateMessage, request, nil)
	resp, err := e.SendRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	var result types.CreateMessageResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (e *BaseMcpEndpoint) RequestSamplingWithChatMessage(ctx context.Context, messages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	samplingMessages := []types.SamplingMessage{}
	var systemPrompt *strings.Builder
	for _, message := range messages {
		if message.Role == chatcompletion.RoleSystem {
			if systemPrompt == nil {
				systemPrompt = &strings.Builder{}
			} else {
				systemPrompt.WriteString("\n")
			}

			systemPrompt.WriteString(message.Text())
			continue
		}

		if message.Role == chatcompletion.RoleUser || message.Role == chatcompletion.RoleAssistant {
			role := types.RoleUser
			if message.Role == chatcompletion.RoleAssistant {
				role = types.RoleAssistant
			}
			for _, content := range message.Contents {
				switch con := content.(type) {
				case *contents.TextContent:
					samplingMessages = append(samplingMessages, types.SamplingMessage{
						Content: types.Content{
							Type: "text",
							Text: &con.Text,
						},
						Role: role,
					})
				case *contents.DataContent:
					if con.MediaTypeStartsWith("image") || con.MediaTypeStartsWith("audio") {
						t := "image"
						if con.MediaTypeStartsWith("audio") {
							t = "audio"
						}
						decoded := base64.URLEncoding.EncodeToString(con.Data)
						samplingMessages = append(samplingMessages, types.SamplingMessage{
							Content: types.Content{
								Type:     t,
								MimeType: &con.MediaType,
								Data:     &decoded,
							},
							Role: role,
						})
					}
				}
			}
		}
	}

	var modelPreferences types.ModelPreferences
	if options != nil && options.ModelId != nil {
		modelPreferences = types.ModelPreferences{
			Hints: []types.ModelHint{{
				Name: options.ModelId,
			}},
		}
	}

	systemPromptString := systemPrompt.String()
	request := types.CreateMessageRequestParams{
		RequestParams:    types.RequestParams{},
		MaxTokens:        options.MaxOutputTokens,
		Messages:         samplingMessages,
		Metadata:         nil,
		ModelPreferences: modelPreferences,
		StopSequences:    options.StopSequences,
		SystemPrompt:     &systemPromptString,
		Temperature:      options.Temperature,
	}

	result, err := e.RequestSampling(ctx, request)
	if err != nil {
		return nil, err
	}

	message := &chatcompletion.ChatMessage{
		Contents: []contents.IAIContent{mcp.ContentToAIContent(result.Content)},
	}
	if result.Role == types.RoleUser {
		message.Role = chatcompletion.RoleUser
	}
	if result.Role == types.RoleAssistant {
		message.Role = chatcompletion.RoleAssistant
	}
	resp := chatcompletion.NewChatResponse(nil, message)
	if result.StopReason != nil {
		if *result.StopReason == "maxTokens" {
			t := chatcompletion.ReasonLength
			resp.FinishReason = &t
		} else {

			t := chatcompletion.ReasonStop
			resp.FinishReason = &t
		}
	}
	return resp, nil
}
