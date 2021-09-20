package main

import "fmt"

func main() {
	arr := []int{9, 3, 4, 5, 7, 1, 33, 69, 94, 84, 67}
	result := quicksort(arr)
	fmt.Println(result)
}

func quicksort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	piovt := arr[0]
	var left, right []int
	for _, ele := range arr[1:] {
		if ele <= piovt {
			left = append(left, ele)
		} else {
			right = append(right, ele)
		}
	}
	return append(quicksort(left), append([]int{piovt}, quicksort(right)...)...)
}
