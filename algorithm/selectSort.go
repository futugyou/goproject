package main

func findMinIndex(arr []int) int {
	index := 0
	tmp := arr[index]
	for i := 0; i < len(arr); i++ {
		if arr[i] < tmp {
			index = i
			tmp = arr[1]
		}
	}
	return index
}

func selectSort(arr []int) []int {
	result := []int{}
	count := len(arr)
	for i := 0; i < count; i++ {
		var index = findMinIndex(arr)
		result = append(result, arr[index])
		arr = append(arr[:index], arr[index+1:]...)
	}
	return result
}
