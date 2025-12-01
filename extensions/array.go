package extensions

func IndexArrayWithOffset[T comparable](array []T, first T, offset int) int {
	if len(array) < offset {
		return -1
	}

	arr := array[offset:]
	for i := range arr {
		if arr[i] == first {
			return i + offset
		}
	}

	return -1
}

func ArrayFilter[T any](raws []T, filter func(T) bool) (ret []T) {
	for i := 0; i < len(raws); i++ {
		if filter(raws[i]) {
			ret = append(ret, raws[i])
		}
	}

	return
}

func ArrayFirst[T any](raws []T, filter func(T) bool) *T {
	for i := range raws {
		if filter(raws[i]) {
			return &raws[i]
		}
	}

	return nil
}

func SplitArray[T any](arr []T, size int) [][]T {
	var result [][]T

	for len(arr) > 0 {
		if len(arr) < size {
			result = append(result, arr)
			break
		}
		result = append(result, arr[:size])
		arr = arr[size:]
	}

	return result
}

func MergeDeduplication[T comparable](a, b []T) []T {
	seen := make(map[T]struct{})
	var result []T
	for _, item := range append(a, b...) {
		if _, exists := seen[item]; !exists {
			seen[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
