package extensions

func GetMapKeys[Key comparable, Value any](m map[Key]Value) []Key {
	keys := make([]Key, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}

func GetMapValues[Key comparable, Value any](m map[Key]Value) []Value {
	values := make([]Value, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

// wo golang verison is 1.20, can not use slices lib.
func MapToSlice[K comparable, V any, R any](m map[K]V, transform func(K, V) R) []R {
	result := make([]R, 0, len(m))
	for k, v := range m {
		if transform != nil {
			result = append(result, transform(k, v))
		} else {
			result = append(result, any(v).(R))
		}
	}
	return result
}

func SliceToMap[T any, K comparable](slice []T, keySelector func(T) K) map[K]T {
	result := make(map[K]T)
	for _, v := range slice {
		key := keySelector(v)
		result[key] = v
	}
	return result
}

func SliceToMapWithTransform[T any, K comparable, V any](slice []T, keySelector func(T) K, valueTransform func(T) V) map[K]V {
	result := make(map[K]V)
	for _, v := range slice {
		key := keySelector(v)
		value := valueTransform(v)
		result[key] = value
	}
	return result
}

func SliceToMapWithTransform2[T any, K comparable, V any](slice []T, keySelector func(T) K, valueTransform func(T) (*V, bool)) map[K]V {
	result := make(map[K]V)
	for _, v := range slice {
		key := keySelector(v)
		if value, ok := valueTransform(v); ok {
			result[key] = *value
		}
	}
	return result
}
