package code0378

import "log"

func Exection() {
	k := 8
	matrix := [][]int{{1, 5, 9}, {10, 11, 13}, {12, 13, 15}}
	n := len(matrix)
	arr := make([]int, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			arr[i*n+j] = matrix[i][j]
		}
	}
	log.Println(exection(arr, k))
}

func exection(arr []int, k int) int {
	lo := 0
	hi := len(arr) - 1
	for {
		if lo > hi {
			break
		}
		p := partition(arr, lo, hi)
		if p < k-1 {
			lo = p + 1
		} else if p > k-1 {
			hi = p - 1
		} else {
			return arr[p]
		}
	}
	return -1
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
