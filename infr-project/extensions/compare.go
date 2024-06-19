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
