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
