package server

import (
	"sync"
)

// McpServerPrimitiveCollection is a thread-safe collection of primitives indexed by name
type McpServerPrimitiveCollection[T IMcpServerPrimitive] struct {
	primitives sync.Map // map[string]T
	mu         sync.RWMutex
	changed    []func()
}

// NewMcpServerPrimitiveCollection creates a new instance of McpServerPrimitiveCollection
func NewMcpServerPrimitiveCollection[T IMcpServerPrimitive]() *McpServerPrimitiveCollection[T] {
	return &McpServerPrimitiveCollection[T]{
		primitives: sync.Map{},
	}
}

// Count returns the number of primitives in the collection
func (c *McpServerPrimitiveCollection[T]) Count() int {
	count := 0
	c.primitives.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// IsEmpty returns true if the collection contains no primitives
func (c *McpServerPrimitiveCollection[T]) IsEmpty() bool {
	count := 0
	c.primitives.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count == 0
}

// OnChanged registers a callback that will be invoked when the collection changes
func (c *McpServerPrimitiveCollection[T]) OnChanged(callback func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.changed = append(c.changed, callback)
}

// raiseChanged invokes all registered change callbacks
func (c *McpServerPrimitiveCollection[T]) raiseChanged() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, callback := range c.changed {
		callback()
	}
}

// Get returns the primitive with the specified name
func (c *McpServerPrimitiveCollection[T]) Get(name string) (T, bool) {
	if name == "" {
		var zero T
		return zero, false
	}

	value, ok := c.primitives.Load(name)
	if !ok {
		var zero T
		return zero, false
	}
	return value.(T), true
}

// Clear removes all primitives from the collection
func (c *McpServerPrimitiveCollection[T]) Clear() {
	c.primitives.Range(func(key, _ interface{}) bool {
		c.primitives.Delete(key)
		return true
	})
	c.raiseChanged()
}

// Add adds a primitive to the collection
func (c *McpServerPrimitiveCollection[T]) Add(primitive T) error {
	if _, ok := c.Get(primitive.GetId()); ok {
		return &AlreadyExistsError{Name: primitive.GetId()}
	}

	c.primitives.Store(primitive.GetId(), primitive)
	c.raiseChanged()
	return nil
}

// TryAdd attempts to add a primitive to the collection
func (c *McpServerPrimitiveCollection[T]) TryAdd(primitive T) bool {
	if _, ok := c.Get(primitive.GetId()); ok {
		return false
	}

	c.primitives.Store(primitive.GetId(), primitive)
	c.raiseChanged()
	return true
}

// Remove removes a primitive from the collection
func (c *McpServerPrimitiveCollection[T]) Remove(primitive T) bool {
	existing, ok := c.Get(primitive.GetId())
	if !ok || !isEqual(existing, primitive) {
		return false
	}

	c.primitives.Delete(primitive.GetId())
	c.raiseChanged()
	return true
}

// Contains checks if a primitive exists in the collection
func (c *McpServerPrimitiveCollection[T]) Contains(primitive T) bool {
	existing, ok := c.Get(primitive.GetId())
	return ok && isEqual(existing, primitive)
}

// Names returns all primitive names in the collection
func (c *McpServerPrimitiveCollection[T]) Names() []string {
	var names []string
	c.primitives.Range(func(key, _ interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

// ToSlice returns all primitives as a slice
func (c *McpServerPrimitiveCollection[T]) ToSlice() []T {
	var slice []T
	c.primitives.Range(func(_, value interface{}) bool {
		slice = append(slice, value.(T))
		return true
	})
	return slice
}

// AlreadyExistsError is returned when trying to add a duplicate primitive
type AlreadyExistsError struct {
	Name string
}

func (e *AlreadyExistsError) Error() string {
	return "a primitive with the same name '" + e.Name + "' already exists in the collection"
}

// Helper function to compare primitives (implementation depends on your needs)
func isEqual[T IMcpServerPrimitive](a, b T) bool {
	return a.GetId() == b.GetId()
}
