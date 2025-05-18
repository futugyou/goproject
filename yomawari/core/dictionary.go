package core

type Dictionary[K comparable, V any] struct {
	data map[K]V
}

func NewDictionary[K comparable, V any]() *Dictionary[K, V] {
	return &Dictionary[K, V]{
		data: make(map[K]V),
	}
}

// Get returns the value associated with the key
func (a *Dictionary[K, V]) Get(key K) (V, bool) {
	val, ok := a.data[key]
	return val, ok
}

func (a *Dictionary[K, V]) Items() map[K]V {
	return a.data
}

// Set sets the value for a key
func (a *Dictionary[K, V]) Set(key K, value V) {
	a.data[key] = value
}

// ContainsKey checks if the key exists in data
func (a *Dictionary[K, V]) ContainsKey(key K) bool {
	_, ok := a.data[key]
	return ok
}

// Remove removes a key from data
func (a *Dictionary[K, V]) Remove(key K) {
	delete(a.data, key)
}

// Clear removes all data
func (a *Dictionary[K, V]) Clear() {
	a.data = make(map[K]V)
}

// Keys returns all keys in data
func (a *Dictionary[K, V]) Keys() []K {
	keys := make([]K, 0, len(a.data))
	for k := range a.data {
		keys = append(keys, k)
	}
	return keys
}

// Values returns all values in data
func (a *Dictionary[K, V]) Values() []V {
	values := make([]V, 0, len(a.data))
	for _, v := range a.data {
		values = append(values, v)
	}
	return values
}

// Count returns the number of data
func (a *Dictionary[K, V]) Count() int {
	return len(a.data)
}
