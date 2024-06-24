package extensions

func StringArrayCompare(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	counts := make(map[string]int)

	for _, s := range arr1 {
		counts[s]++
	}

	for _, s := range arr2 {
		if counts[s] == 0 {
			return false
		}
		counts[s]--
	}

	for _, count := range counts {
		if count != 0 {
			return false
		}
	}

	return true
}

func MapsCompare(map1, map2 map[string]string) bool {
	if len(map1) != len(map2) {
		return false
	}

	for key, value1 := range map1 {
		value2, ok := map2[key]
		if !ok || value1 != value2 {
			return false
		}
	}

	return true
}
