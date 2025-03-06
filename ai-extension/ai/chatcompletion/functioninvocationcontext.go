package chatcompletion

import (
	"github.com/futugyou/ai-extension/abstractions/chatcompletion"
	"github.com/futugyou/ai-extension/abstractions/contents"
	"github.com/futugyou/ai-extension/abstractions/functions"
)

type FunctionInvocationContext struct {
	Messages          []chatcompletion.ChatMessage
	CallContent       contents.FunctionCallContent
	Options           *chatcompletion.ChatOptions
	Function          functions.AIFunction
	Iteration         int
	FunctionCallIndex int
	FunctionCount     int
	Terminate         bool
}

func NewFunctionInvocationContext() *FunctionInvocationContext {
	return &FunctionInvocationContext{
		Messages: make([]chatcompletion.ChatMessage, 0),
	}
}
