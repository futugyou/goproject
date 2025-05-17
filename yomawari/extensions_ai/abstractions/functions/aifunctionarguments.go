package functions

import "github.com/futugyou/yomawari/core"

// AIFunctionArguments represents arguments for AI functions with additional context capabilities
type AIFunctionArguments struct {
	core.Dictionary[string, any]

	// Additional context associated with these arguments
	Context map[interface{}]interface{}
}

// NewAIFunctionArguments creates a new empty AIFunctionArguments instance
func NewAIFunctionArguments() *AIFunctionArguments {
	return &AIFunctionArguments{
		Dictionary: *core.NewDictionary[string, any](),
		Context:    make(map[interface{}]interface{}),
	}
}

// NewAIFunctionArgumentsFromMap creates a new AIFunctionArguments instance from an existing map
func NewAIFunctionArgumentsFromMap(arguments map[string]interface{}) *AIFunctionArguments {
	args := NewAIFunctionArguments()
	for k, v := range arguments {
		args.Set(k, v)
	}

	return args
}

// GetContext returns a value from context by key
func (a *AIFunctionArguments) GetContext(key interface{}) (interface{}, bool) {
	val, ok := a.Context[key]
	return val, ok
}

// SetContext sets a value in context
func (a *AIFunctionArguments) SetContext(key, value interface{}) {
	if a.Context == nil {
		a.Context = make(map[interface{}]interface{})
	}
	a.Context[key] = value
}

// RemoveContext removes a key from context
func (a *AIFunctionArguments) RemoveContext(key interface{}) {
	delete(a.Context, key)
}

// ClearContext clears all context
func (a *AIFunctionArguments) ClearContext() {
	a.Context = make(map[interface{}]interface{})
}
