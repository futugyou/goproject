package chatcompletion

import (
	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/functions"
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
