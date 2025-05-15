package functions

type KernelArguments struct {
	arguments map[string]interface{}
}

// Get returns the value associated with the key
func (a *KernelArguments) Get(key string) (interface{}, bool) {
	val, ok := a.arguments[key]
	return val, ok
}

func (a *KernelArguments) Gets() map[string]interface{} {
	return a.arguments
}

// Set sets the value for a key
func (a *KernelArguments) Set(key string, value interface{}) {
	a.arguments[key] = value
}

// ContainsKey checks if the key exists in arguments
func (a *KernelArguments) ContainsKey(key string) bool {
	_, ok := a.arguments[key]
	return ok
}

// Remove removes a key from arguments
func (a *KernelArguments) Remove(key string) {
	delete(a.arguments, key)
}

// Clear removes all arguments
func (a *KernelArguments) Clear() {
	a.arguments = make(map[string]interface{})
}

// Keys returns all keys in arguments
func (a *KernelArguments) Keys() []string {
	keys := make([]string, 0, len(a.arguments))
	for k := range a.arguments {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values in arguments
func (a *KernelArguments) Values() []interface{} {
	values := make([]interface{}, 0, len(a.arguments))
	for _, v := range a.arguments {
		values = append(values, v)
	}
	return values
}

// Count returns the number of arguments
func (a *KernelArguments) Count() int {
	return len(a.arguments)
}
