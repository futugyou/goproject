package code0912

import "fmt"

func Exection() {
	arr := []int{9, 0, 3, 5, 12, -1}
	exection(arr, 0, len(arr)-1)
	fmt.Println(arr)
}

func exection(arr []int, lo, hi int) {
	if lo > hi {
		return
	}
	p := partotion(arr, lo, hi)
	exection(arr, lo, p-1)
	exection(arr, p+1, hi)
}

func partotion(arr []int, lo, hi int) int {
	pivot := arr[lo]
	i := lo + 1
	j := hi
	for {
		for ; i < hi && arr[i] <= pivot; i++ {
		}
		for ; j > lo && arr[j] > pivot; j-- {

		}
		if i >= j {
			break
		}
		swap(arr, i, j)
	}
	swap(arr, lo, j)
	return j
}

func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
