package code0704

import "fmt"

func Exection() {
	arr := []int{-1, 0, 3, 5, 9, 12}
	target := 9
	exection(arr, target)
}

func exection(arr []int, target int) {
	left := 0
	right := len(arr) - 1
	for {
		mid := left + (right-left)/2
		if arr[mid] == target {
			fmt.Println(mid)
			break
		} else if arr[mid] > target {
			right = right - 1
		} else if arr[mid] < target {
			left = left + 1
		}
		if left > right {
			fmt.Println(-1)
			break
		}
	}
}
