package abstractions

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/futugyou/yomawari/semantic_kernel/utilities"
)

type FunctionCallContentBuilder struct {
	functionCallIdsByIndex          map[string]string
	functionNamesByIndex            map[string]string
	functionArgumentBuildersByIndex map[string]*strings.Builder
}

func (b *FunctionCallContentBuilder) Append(content StreamingChatMessageContent) {
	for _, item := range content.Items.Items() {
		if con, ok := item.(StreamingFunctionCallUpdateContent); ok && item.Type() == "streaming-function-call-update" {
			trackStreamingFunctionCallUpdate(con, b.functionCallIdsByIndex, b.functionNamesByIndex, b.functionArgumentBuildersByIndex)
		}
	}
}

func trackStreamingFunctionCallUpdate(update StreamingFunctionCallUpdateContent, functionCallIdsByIndex map[string]string,
	functionNamesByIndex map[string]string, functionArgumentBuildersByIndex map[string]*strings.Builder) {

	var functionCallIndex = fmt.Sprintf("%d-%d", update.RequestIndex, update.FunctionCallIndex)
	if len(update.CallId) > 0 {
		functionCallIdsByIndex[functionCallIndex] = update.CallId
	}

	// Ensure we're tracking the function's name.
	if len(update.Name) > 0 {
		functionNamesByIndex[functionCallIndex] = update.Name
	}

	// Ensure we're tracking the function's arguments.
	if len(update.Arguments) > 0 {
		var arguments *strings.Builder
		var ok bool
		if arguments, ok = functionArgumentBuildersByIndex[functionCallIndex]; !ok {
			arguments = &strings.Builder{}
			functionArgumentBuildersByIndex[functionCallIndex] = arguments
		}

		arguments.WriteString(update.Arguments)
	}
}

func (b *FunctionCallContentBuilder) getFunctionArguments(functionCallIndex string) (map[string]any, error) {
	var functionArgumentsBuilder *strings.Builder
	var ok bool
	if functionArgumentsBuilder, ok = b.functionArgumentBuildersByIndex[functionCallIndex]; !ok {
		return nil, nil
	}

	argumentsString := functionArgumentsBuilder.String()
	if len(argumentsString) == 0 {
		return nil, nil
	}

	var arguments map[string]any
	if err := json.Unmarshal([]byte(argumentsString), &arguments); err != nil {
		return nil, err
	}

	return arguments, nil
}
func (b *FunctionCallContentBuilder) Build() []FunctionCallContent {
	functionCalls := []FunctionCallContent{}
	for key, functionCallIndexAndId := range b.functionCallIdsByIndex {
		pluginName := ""
		functionName := ""

		if fqn, ok := b.functionNamesByIndex[key]; ok {
			functionFullyQualifiedName := utilities.ParseFunctionName(fqn, "")
			pluginName = functionFullyQualifiedName.PluginName
			functionName = functionFullyQualifiedName.Name
		}

		arguments, err := b.getFunctionArguments(key)

		ff := FunctionCallContent{
			MimeType:     "",
			ModelId:      "",
			Metadata:     map[string]any{},
			Id:           functionCallIndexAndId,
			PluginName:   pluginName,
			FunctionName: functionName,
			Arguments:    arguments,
			Exception:    "",
			InnerContent: nil,
		}

		if err != nil {
			ff.Exception = err.Error()
		}

		functionCalls = append(functionCalls, ff)
	}

	return functionCalls
}
