package code0215

import "fmt"

func Exection() {
	arr := []int{2, 1, 5, 4}
	k := 1
	exection(arr, k)
}

func exection(arr []int, k int) {
	k = len(arr) - k
	lo := 0
	hi := len(arr) - 1
	for {
		if lo > hi {
			break
		}
		p := partition(arr, lo, hi)
		if p < k {
			lo = p + 1
		} else if p > k {
			hi = p - 1
		} else {
			fmt.Println(arr[p])
			return
		}
	}
	fmt.Println(-1)
}

func partition(arr []int, lo, hi int) int {
	p := arr[lo]
	i := lo + 1
	j := hi
	for {
		if i > j {
			break
		}
		for ; i < hi && arr[i] < p; i++ {
		}
		for ; j > lo && arr[j] > p; j-- {
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
