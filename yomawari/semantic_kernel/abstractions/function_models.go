package abstractions

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/futugyou/yomawari/core"
)

type KernelArguments struct {
	core.Dictionary[string, any]
}

type FunctionResult struct {
	Function       KernelFunction
	Metadata       map[string]any
	RenderedPrompt string
	Value          any
}

func (m *FunctionResult) ValueType() (reflect.Type, string) {
	if m.Value == nil {
		return nil, ""
	}
	r := reflect.TypeOf(m.Value)
	return r, r.String()
}

func GetFunctionResultValue[T any](r FunctionResult) (T, error) {
	var zero T

	if r.Value == nil {
		return zero, nil
	}

	// Case 1: Value is T
	if v, ok := r.Value.(T); ok {
		return v, nil
	}

	// Case 2: Value is KernelContent
	if content, ok := r.Value.(KernelContent); ok {
		tType := reflect.TypeOf(zero)
		if tType.Kind() == reflect.String {
			return any(content.ToString()).(T), nil
		}

		if inner := content.GetInnerContent(); inner != nil {
			if ic, ok := inner.(T); ok {
				return ic, nil
			}
		}
	}

	return zero, fmt.Errorf("cannot cast %T to %T", r.Value, zero)
}

var FunctionResultTypeRegistry map[string]reflect.Type

// data := []byte(`123`)
//
// v, err := Deserialize("int", data)
//
// fmt.Printf("value=%v, type=%T\n", v, v)
func Deserialize(typeName string, data json.RawMessage) (interface{}, error) {
	typ, ok := FunctionResultTypeRegistry[typeName]
	if !ok {
		return nil, fmt.Errorf("unknown type: %s", typeName)
	}

	ptr := reflect.New(typ).Interface()

	err := json.Unmarshal(data, ptr)
	if err != nil {
		return nil, err
	}

	return reflect.ValueOf(ptr).Elem().Interface(), nil
}

type KernelFunctionSchemaModel struct {
	Type        string                            `json:"type"`
	Description string                            `json:"condition"`
	Properties  map[string]map[string]interface{} `json:"properties"`
	Required    []string                          `json:"required"`
}
