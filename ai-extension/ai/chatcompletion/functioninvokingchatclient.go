package chatcompletion

import (
	"context"
	"fmt"
	"sync"

	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
	"github.com/futugyou/ai-extension/abstractions/contents"
	"github.com/futugyou/ai-extension/abstractions/functions"
)

type FunctionInvokingChatClient struct {
	chatcompletion.DelegatingChatClient
	MaximumIterationsPerRequest int
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
		return newFunctionInvocationResult("Continue", "NotFound", callContent, nil, nil)
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
		return newFunctionInvocationResult(continueMode, "Exception", callContent, nil, nil)
	}

	continueMode := "Terminate"
	if !functionContext.Terminate {
		continueMode = "Continue"
	}
	return newFunctionInvocationResult(continueMode, "RanToCompletion", callContent, result, nil)
}

func (f *FunctionInvokingChatClient) InvokeFunction(ctx context.Context, funcContext FunctionInvocationContext) interface{} {
	f.CurrentContext = &funcContext
	if funcContext.Function == nil {
		return nil
	}

	resullt, err := funcContext.Function.InvokeAsync(ctx, funcContext.CallContent.Arguments)
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
			AIContent: contents.AIContent{},
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

func newFunctionInvocationResult(
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
