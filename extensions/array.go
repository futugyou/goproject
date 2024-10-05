package extensions

func IndexArrayWithOffset[T comparable](array []T, first T, offset int) int {
	if len(array) < offset {
		return -1
	}

	arr := array[offset:]
	for i := 0; i < len(arr); i++ {
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
