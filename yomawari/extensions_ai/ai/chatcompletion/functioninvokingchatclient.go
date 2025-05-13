package chatcompletion

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
	"github.com/google/uuid"
)

type FunctionInvokingChatClient struct {
	chatcompletion.DelegatingChatClient
	MaximumIterationsPerRequest *int
	RetryOnError                bool
	IncludeDetailedErrors       bool
	AllowConcurrentInvocation   bool
	CurrentContext              *FunctionInvocationContext
}

func NewFunctionInvokingChatClient(innerClient chatcompletion.IChatClient) *FunctionInvokingChatClient {
	return &FunctionInvokingChatClient{
		DelegatingChatClient: chatcompletion.DelegatingChatClient{InnerClient: innerClient},
	}
}

func (f *FunctionInvokingChatClient) processFunctionCalls(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
	callContents []contents.FunctionCallContent,
	iteration int,
) (string, []chatcompletion.ChatMessage, []chatcompletion.ChatMessage, error) {
	if len(messages) == 0 {
		return "", nil, messages, fmt.Errorf("no messages to process")
	}

	if len(messages) == 1 {
		result := f.processFunctionCall(ctx, messages, options, callContents, iteration, 0)
		added, err := f.CreateResponseMessages([]FunctionInvocationResult{result})
		if err != nil {
			return "", nil, messages, err
		}

		messages = append(messages, added...)

		return result.ContinueMode, added, messages, nil
	}

	var results []FunctionInvocationResult
	if f.AllowConcurrentInvocation {
		var wg sync.WaitGroup
		resultChan := make(chan FunctionInvocationResult, len(callContents))

		for i := 0; i < len(callContents); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				res := f.processFunctionCall(ctx, messages, options, callContents, iteration, i)
				resultChan <- res
			}(i)
		}

		wg.Wait()
		close(resultChan)

		results = make([]FunctionInvocationResult, 0, len(callContents))
		for res := range resultChan {
			results = append(results, res)
		}
	} else {
		results = make([]FunctionInvocationResult, len(callContents))
		for i := 0; i < len(callContents); i++ {
			results[i] = f.processFunctionCall(ctx, messages, options, callContents, iteration, i)
		}
	}

	added, err := f.CreateResponseMessages(results)
	if err != nil {
		return "", nil, messages, err
	}

	continueMode := "Continue"
	messages = append(messages, added...)
	for _, fir := range results {
		if fir.ContinueMode == "Terminate" {
			continueMode = "Terminate"
			break
		} else if fir.ContinueMode == "AllowOneMoreRoundtrip" {
			continueMode = "AllowOneMoreRoundtrip"
		}
	}

	return continueMode, added, messages, nil
}

func (f *FunctionInvokingChatClient) processFunctionCall(
	ctx context.Context,
	messages []chatcompletion.ChatMessage,
	options *chatcompletion.ChatOptions,
	callContents []contents.FunctionCallContent,
	iteration int,
	functionCallIndex int,
) FunctionInvocationResult {
	var callContent = callContents[functionCallIndex]
	var function functions.AIFunction
	if options != nil && len(options.Tools) > 0 {
		for i := 0; i < len(options.Tools); i++ {
			if t, ok := options.Tools[i].(functions.AIFunction); ok {
				if t.GetName() == callContent.Name {
					function = t
					break
				}
			}
		}
	}

	if function == nil {
		return NewFunctionInvocationResult("Continue", "NotFound", callContent, nil, nil)
	}

	functionContext := FunctionInvocationContext{
		Messages:          messages,
		CallContent:       callContent,
		Options:           options,
		Function:          function,
		Iteration:         iteration,
		FunctionCallIndex: functionCallIndex,
	}

	result := f.InvokeFunction(ctx, functionContext)
	if result == nil {
		continueMode := "Continue"
		if !f.RetryOnError {
			continueMode = "AllowOneMoreRoundtrip"
		}
		return NewFunctionInvocationResult(continueMode, "Exception", callContent, nil, nil)
	}

	continueMode := "Terminate"
	if !functionContext.Terminate {
		continueMode = "Continue"
	}
	return NewFunctionInvocationResult(continueMode, "RanToCompletion", callContent, result, nil)
}

func (f *FunctionInvokingChatClient) InvokeFunction(ctx context.Context, funcContext FunctionInvocationContext) interface{} {
	f.CurrentContext = &funcContext
	if funcContext.Function == nil {
		return nil
	}
	arguments := functions.NewAIFunctionArgumentsFromMap(funcContext.CallContent.Arguments)
	resullt, err := funcContext.Function.Invoke(ctx, *arguments)
	if err != nil {
		// log error
		return nil
	}

	return resullt
}

func (f *FunctionInvokingChatClient) CreateResponseMessages(results []FunctionInvocationResult) ([]chatcompletion.ChatMessage, error) {
	conts := []contents.IAIContent{}
	createFunctionResultContent := func(result FunctionInvocationResult) contents.FunctionResultContent {
		var functionResult interface{}
		if result.Status == "RanToCompletion" {
			functionResult = result.Result
			if functionResult == nil {
				functionResult = "Success: Function completed."
			}
		} else {
			functionResult = fmt.Sprintf("Error: %s", result.Status)
		}

		return contents.FunctionResultContent{
			AIContent: contents.NewAIContent(nil, nil),
			CallId:    result.CallContent.CallId,
			Result:    functionResult,
			Error:     result.err,
		}
	}

	for i := 0; i < len(results); i++ {
		conts = append(conts, createFunctionResultContent(results[i]))
	}
	return []chatcompletion.ChatMessage{
		{
			Role:     chatcompletion.RoleTool,
			Contents: conts,
		},
	}, nil
}

type FunctionInvocationResult struct {
	Status       string
	CallContent  contents.FunctionCallContent
	Result       interface{}
	err          error
	ContinueMode string
}

func NewFunctionInvocationResult(
	continueMode string,
	status string,
	callContent contents.FunctionCallContent,
	result interface{},
	err error,
) FunctionInvocationResult {
	return FunctionInvocationResult{
		ContinueMode: continueMode,
		Status:       status,
		CallContent:  callContent,
		Result:       result,
		err:          err,
	}
}

func (f *FunctionInvokingChatClient) updateOptionsForMode(mode string, options *chatcompletion.ChatOptions, chatThreadId *string) bool {
	switch mode {
	case "Continue":
		if *options.ToolMode == chatcompletion.RequireAnyMode {
			options.ToolMode = nil
			options.ChatThreadId = chatThreadId
		}
		return false
	case "Terminate":
		return true
	case "AllowOneMoreRoundtrip":
		options.Tools = nil
		options.ToolMode = nil
		options.ChatThreadId = chatThreadId
		return false
	default:
		if options.ChatThreadId != chatThreadId {
			options.ChatThreadId = chatThreadId
		}
	}

	return false
}

func (f *FunctionInvokingChatClient) copyFunctionCalls(messages []chatcompletion.ChatMessage, functionCalls *[]contents.FunctionCallContent) bool {
	any := false
	for _, message := range messages {
		if f.copyFunctionCallsFromContent(message.Contents, functionCalls) {
			any = true
		}
	}
	return any
}

func (f *FunctionInvokingChatClient) copyFunctionCallsFromContent(content []contents.IAIContent, functionCalls *[]contents.FunctionCallContent) bool {
	any := false
	for _, item := range content {
		if functionCall, ok := item.(contents.FunctionCallContent); ok {
			*functionCalls = append(*functionCalls, functionCall)
			any = true
		}
	}
	return any
}

func (f *FunctionInvokingChatClient) fixupHistories(
	originalMessages []chatcompletion.ChatMessage,
	augmentedHistory []chatcompletion.ChatMessage,
	response chatcompletion.ChatResponse,
	allTurnsResponseMessages []chatcompletion.ChatMessage,
	lastIterationHadThreadId bool,
) ([]chatcompletion.ChatMessage, []chatcompletion.ChatMessage, bool) {
	if response.ChatThreadId != nil {
		augmentedHistory = augmentedHistory[:0]
		lastIterationHadThreadId = true
	} else if lastIterationHadThreadId {
		augmentedHistory = append([]chatcompletion.ChatMessage{}, originalMessages...)
		augmentedHistory = append(augmentedHistory, allTurnsResponseMessages...)
		lastIterationHadThreadId = false
	} else {
		if len(augmentedHistory) == 0 {
			augmentedHistory = append([]chatcompletion.ChatMessage{}, originalMessages...)
		}

		addMessages := func(messages []chatcompletion.ChatMessage, response chatcompletion.ChatResponse) []chatcompletion.ChatMessage {
			messages = append(messages, response.Messages...)
			return messages
		}

		augmentedHistory = addMessages(augmentedHistory, response)
		lastIterationHadThreadId = false
	}

	return augmentedHistory, augmentedHistory, lastIterationHadThreadId
}

func (c *FunctionInvokingChatClient) GetStreamingResponse(ctx context.Context, messages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) <-chan chatcompletion.ChatStreamingResponse {
	updateCh := make(chan chatcompletion.ChatStreamingResponse)
	if messages == nil {
		updateCh <- chatcompletion.ChatStreamingResponse{
			Update: nil,
			Err:    errors.New("messages cannot be nil"),
		}
		close(updateCh)
		return updateCh
	}

	go func() {
		defer close(updateCh)

		var (
			originalMessages    = messages
			augmentedHistory    []chatcompletion.ChatMessage
			functionCallContent []contents.FunctionCallContent
			responseMessages    []chatcompletion.ChatMessage
			lastIterationHadID  bool
		)

		for iteration := 0; ; iteration++ {
			var updates []chatcompletion.ChatResponseUpdate
			functionCallContent = nil

			innerCh := c.InnerClient.GetStreamingResponse(ctx, messages, options)

			for update := range innerCh {
				if update.Err != nil {
					updateCh <- chatcompletion.ChatStreamingResponse{
						Update: nil,
						Err:    update.Err,
					}
					return
				}

				select {
				case <-ctx.Done():
					return
				case updateCh <- update:

				}

				updates = append(updates, *update.Update)
				c.copyFunctionCallsFromContent(update.Update.Contents, &functionCallContent)
			}

			if len(functionCallContent) == 0 || options == nil || len(options.Tools) == 0 || (c.MaximumIterationsPerRequest != nil && iteration >= *c.MaximumIterationsPerRequest) {
				return
			}

			response := chatcompletion.ToChatResponse(updates)
			responseMessages = append(responseMessages, response.Messages...)
			messages, augmentedHistory, lastIterationHadID = c.fixupHistories(originalMessages, augmentedHistory, response, responseMessages, lastIterationHadID)

			continueMode := ""
			var modeAndMessages []chatcompletion.ChatMessage
			var err error
			continueMode, modeAndMessages, augmentedHistory, err = c.processFunctionCalls(ctx, augmentedHistory, options, functionCallContent, iteration)
			if err != nil {
				updateCh <- chatcompletion.ChatStreamingResponse{
					Update: nil,
					Err:    err,
				}
				return
			}

			responseMessages = append(responseMessages, modeAndMessages...)

			toolResponseID := uuid.New().String()
			for _, msg := range modeAndMessages {
				t := time.Now().UTC()
				toolResultUpdate := chatcompletion.ChatResponseUpdate{
					AdditionalProperties: msg.AdditionalProperties,
					AuthorName:           msg.AuthorName,
					ChatThreadId:         response.ChatThreadId,
					CreatedAt:            &t,
					Contents:             msg.Contents,
					ResponseId:           &toolResponseID,
					Role:                 &msg.Role,
				}

				select {
				case <-ctx.Done():
					return
				case updateCh <- chatcompletion.ChatStreamingResponse{
					Update: &toolResultUpdate,
					Err:    nil,
				}:
				}
			}

			if c.updateOptionsForMode(continueMode, options, response.ChatThreadId) {
				return
			}
		}
	}()

	return updateCh
}

func (c *FunctionInvokingChatClient) GetResponse(ctx context.Context, chatMessages []chatcompletion.ChatMessage, options *chatcompletion.ChatOptions) (*chatcompletion.ChatResponse, error) {
	if chatMessages == nil {
		return nil, errors.New("messages cannot be nil")
	}

	var originalMessages = chatMessages
	var augmentedHistory []chatcompletion.ChatMessage = []chatcompletion.ChatMessage{}
	var functionCallContent []contents.FunctionCallContent = []contents.FunctionCallContent{}
	var responseMessages []chatcompletion.ChatMessage = []chatcompletion.ChatMessage{}
	var lastIterationHadID bool
	var response *chatcompletion.ChatResponse
	var err error
	totalUsage := &abstractions.UsageDetails{}

	for iteration := 0; ; iteration++ {
		functionCallContent = nil

		response, err = c.InnerClient.GetResponse(ctx, chatMessages, options)
		if err != nil {
			return nil, err
		}

		requiresFunctionInvocation := len(options.Tools) > 0 &&
			(c.MaximumIterationsPerRequest == nil || iteration < *c.MaximumIterationsPerRequest) &&
			c.copyFunctionCalls(response.Messages, &functionCallContent)

		if iteration == 0 && !requiresFunctionInvocation {
			return response, nil
		}

		responseMessages = append(responseMessages, response.Messages...)

		if response.Usage != nil {
			totalUsage.AddUsageDetails(*response.Usage)
		}

		if !requiresFunctionInvocation {
			break
		}

		chatMessages, augmentedHistory, lastIterationHadID = c.fixupHistories(originalMessages, augmentedHistory, *response, responseMessages, lastIterationHadID)

		continueMode := ""
		var modeAndMessages []chatcompletion.ChatMessage

		continueMode, modeAndMessages, augmentedHistory, err = c.processFunctionCalls(ctx, augmentedHistory, options, functionCallContent, iteration)
		if err != nil {
			return nil, err
		}

		responseMessages = append(responseMessages, modeAndMessages...)

		if c.updateOptionsForMode(continueMode, options, response.ChatThreadId) {
			break
		}
	}

	response.Messages = responseMessages
	response.Usage = totalUsage

	return response, nil
}
