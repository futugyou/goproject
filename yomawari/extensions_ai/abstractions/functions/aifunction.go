package functions

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/futugyou/yomawari/extensions_ai/abstractions"
)

type AIFunction interface {
	abstractions.AITool
	Invoke(ctx context.Context, arguments AIFunctionArguments) (interface{}, error)
	// i want use https://github.com/invopop/jsonschema
	GetJsonSchema() map[string]interface{}
	UnderlyingMethod() reflect.Method
	JsonSerializerOptions() json.Marshaler
}

type BaseAIFunction struct {
	*abstractions.BaseAITool
}

func (t *BaseAIFunction) GetJsonSchema() map[string]interface{} {
	return map[string]interface{}{}
}

func (f *BaseAIFunction) UnderlyingMethod() reflect.Method {
	return reflect.TypeOf(f).Method(0)
}

func (f *BaseAIFunction) JsonSerializerOptions() json.Marshaler {
	return &json.RawMessage{}
}
