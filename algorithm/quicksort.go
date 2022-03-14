package main

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
func exection(arr []int, target int) {
	panic("unimplemented")
}
