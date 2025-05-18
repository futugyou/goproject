package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

var _ AIFunction = (*AddNumbersFunction)(nil)

type AddNumbersFunction struct {
	BaseAIFunction
}

func (f *AddNumbersFunction) Invoke(ctx context.Context, arguments AIFunctionArguments) (interface{}, error) {
	var a float64
	var b float64
	if aVal, ok := arguments.Items()["a"].(float64); ok {
		a = aVal
	} else {
		return nil, fmt.Errorf("invalid arguments for 'a'")
	}

	if bVal, ok := arguments.Items()["b"].(float64); ok {
		b = bVal
	} else {
		return nil, fmt.Errorf("invalid arguments for 'b'")
	}

	return a + b, nil
}

func (f *AddNumbersFunction) JsonSchema() json.RawMessage {
	schema := map[string]interface{}{
		"title":       "addNumbers",
		"description": "A simple function that adds two numbers together.",
		"type":        "object",
		"properties": map[string]interface{}{
			"a": map[string]interface{}{"type": "number"},
			"b": map[string]interface{}{"type": "number", "default": 1},
		},
		"required": []string{"a"},
	}
	schemaBytes, _ := json.Marshal(schema)
	return schemaBytes
}

func (f *AddNumbersFunction) UnderlyingMethod() reflect.Method {
	method := reflect.TypeOf(f).Method(0)
	return method
}

func (f *AddNumbersFunction) JsonSerializerOptions() json.Marshaler {
	return &json.RawMessage{}
}

func TestAddNumbersFunction(t *testing.T) {
	function := &AddNumbersFunction{}

	tests := []struct {
		name      string
		arguments map[string]interface{}
		expected  float64
		expectErr bool
	}{
		{
			name:      "Valid arguments",
			arguments: map[string]interface{}{"a": 5, "b": 3},
			expected:  8,
			expectErr: false,
		},
		{
			name:      "Missing 'a' argument",
			arguments: map[string]interface{}{"b": 3},
			expected:  0,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arguments := NewAIFunctionArgumentsFromMap(tt.arguments)
			result, err := function.Invoke(context.Background(), *arguments)
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.expectErr && result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}

	// testing JsonSchema
	expectedSchema := `{"title":"addNumbers","description":"A simple function that adds two numbers together.","type":"object","properties":{"a":{"type":"number"},"b":{"type":"number","default":1}},"required":["a"]}`
	schema := string(function.JsonSchema())
	if schema != expectedSchema {
		t.Errorf("expected schema %v, got %v", expectedSchema, schema)
	}

	// testing UnderlyingMethod
	method := function.UnderlyingMethod()
	if method.Name != "Invoke" {
		t.Errorf("expected method name 'Invoke', got '%v'", method.Name)
	}

	// testing JsonSerializerOptions
	options := function.JsonSerializerOptions()
	if options == nil {
		t.Error("expected JsonSerializerOptions, got nil")
	}
}
