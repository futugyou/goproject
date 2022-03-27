package code0870

import (
	"fmt"
	"sort"
)

func Exection() {
	arr1 := []int{11, 2, 7, 15}
	arr2 := []int{1, 10, 4, 11}
	exection(arr1, arr2)
}

func exection(arr1 []int, arr2 []int) {
	n := len(arr2)
	q := make([][]int, n)
	result := make([]int, n)
	for i := 0; i < n; i++ {
		q[i] = []int{arr2[i], i}
	}
	sort.Slice(q, func(i, j int) bool {
		return q[i][0] > q[j][0]
	})
	sort.Ints(arr1)
	left := 0
	right := n - 1
	for i := 0; i < n; i++ {
		c := q[i]
		if arr1[right] > c[0] {
			result[c[1]] = arr1[right]
			right--
		} else {
			result[c[1]] = arr1[left]
			left++
		}
	}
	fmt.Println(result)
}
