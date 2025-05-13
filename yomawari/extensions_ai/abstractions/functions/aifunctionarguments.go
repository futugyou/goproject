package functions

// AIFunctionArguments represents arguments for AI functions with additional context capabilities
type AIFunctionArguments struct {
	// The nominal arguments
	arguments map[string]interface{}

	// Additional context associated with these arguments
	Context map[interface{}]interface{}
}

// NewAIFunctionArguments creates a new empty AIFunctionArguments instance
func NewAIFunctionArguments() *AIFunctionArguments {
	return &AIFunctionArguments{
		arguments: make(map[string]interface{}),
		Context:   make(map[interface{}]interface{}),
	}
}

// NewAIFunctionArgumentsFromMap creates a new AIFunctionArguments instance from an existing map
func NewAIFunctionArgumentsFromMap(arguments map[string]interface{}) *AIFunctionArguments {
	args := NewAIFunctionArguments()
	if arguments != nil {
		args.arguments = arguments
	}
	return args
}

// Get returns the value associated with the key
func (a *AIFunctionArguments) Get(key string) (interface{}, bool) {
	val, ok := a.arguments[key]
	return val, ok
}

func (a *AIFunctionArguments) GetArguments() map[string]interface{} {
	return a.arguments
}

// Set sets the value for a key
func (a *AIFunctionArguments) Set(key string, value interface{}) {
	a.arguments[key] = value
}

// ContainsKey checks if the key exists in arguments
func (a *AIFunctionArguments) ContainsKey(key string) bool {
	_, ok := a.arguments[key]
	return ok
}

// Remove removes a key from arguments
func (a *AIFunctionArguments) Remove(key string) {
	delete(a.arguments, key)
}

// Clear removes all arguments
func (a *AIFunctionArguments) Clear() {
	a.arguments = make(map[string]interface{})
}

// Keys returns all keys in arguments
func (a *AIFunctionArguments) Keys() []string {
	keys := make([]string, 0, len(a.arguments))
	for k := range a.arguments {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values in arguments
func (a *AIFunctionArguments) Values() []interface{} {
	values := make([]interface{}, 0, len(a.arguments))
	for _, v := range a.arguments {
		values = append(values, v)
	}
	return values
}

// Count returns the number of arguments
func (a *AIFunctionArguments) Count() int {
	return len(a.arguments)
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
