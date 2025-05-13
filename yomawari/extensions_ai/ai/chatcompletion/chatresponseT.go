package chatcompletion

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/futugyou/yomawari/extensions_ai/abstractions/chatcompletion"
	"github.com/futugyou/yomawari/extensions_ai/abstractions/contents"
)

type ChatResponse[T any] struct {
	chatcompletion.ChatResponse
	deserializedResult    *T
	hasDeserializedResult bool
	IsWrappedInObject     bool
}

func NewChatResponse[T any](inner chatcompletion.ChatResponse) *ChatResponse[T] {
	return &ChatResponse[T]{
		ChatResponse: inner,
	}
}

func (response *ChatResponse[T]) GetResult() (*T, error) {
	return response.getResultCore()
}

func (response *ChatResponse[T]) TryGetResult() (*T, bool) {
	result, err := response.getResultCore()
	if err != nil {
		return nil, false
	}
	return result, true
}

func (response *ChatResponse[T]) getResultCore() (*T, error) {
	if response.hasDeserializedResult {
		return response.deserializedResult, nil
	}

	var jsonRaw = response.GetResultAsJson()
	if len(jsonRaw) == 0 {
		return nil, fmt.Errorf("ResultDidNotContainJson")
	}

	var inner json.RawMessage
	var exists bool
	if response.IsWrappedInObject {
		var data map[string]json.RawMessage
		err := json.Unmarshal([]byte(jsonRaw), &data)
		if err != nil {
			return nil, fmt.Errorf("ResultDidNotContainJson")
		}

		inner, exists = data["data"]
		if !exists {
			return nil, fmt.Errorf("ResultDidNotContainDataProperty")
		}
	}

	result, err := DeserializeFirstTopLevelObject2[T](inner)
	if err != nil {
		return nil, fmt.Errorf("DeserializationProducedNull")
	}

	response.deserializedResult = result
	response.hasDeserializedResult = true
	return result, nil
}

func (response *ChatResponse[T]) GetResultAsJson() string {
	if len(response.Messages) != 1 {
		return ""
	}
	if len(response.Messages[0].Contents) != 1 {
		return ""
	}
	content := response.Messages[0].Contents[0]

	switch expression := content.(type) {
	case contents.TextContent:
		return expression.Text
	}
	return ""
}

func DeserializeFirstTopLevelObject[T any](jsonStr string) (*T, error) {
	decoder := json.NewDecoder(bytes.NewReader([]byte(jsonStr)))
	var result T
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func DeserializeFirstTopLevelObject2[T any](jsonData json.RawMessage) (*T, error) {
	decoder := json.NewDecoder(bytes.NewReader(jsonData))
	var result T
	if err := decoder.Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
