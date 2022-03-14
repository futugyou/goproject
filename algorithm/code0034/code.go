package code0034

import "fmt"

func Exection() {
	arr := []int{5, 7, 7, 8, 8, 10}
	target := 8
	exection(arr, target)
}

func exection(arr []int, target int) {
	left := left_bound(arr, target)
	right := right_bound(arr, target)
	fmt.Println(left, right)
}

func right_bound(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	for {
		mid := left + (right-left)/2

		if arr[mid] == target {
			left = left + 1
		} else if arr[mid] < target {
			left = left + 1
		} else if arr[mid] > target {
			right = right - 1
		}
		if left > right {
			break
		}
	}
	if right < 0 || arr[right] != target {
		return -1
	}
	return right
}

func left_bound(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	for {
		mid := left + (right-left)/2
		if arr[mid] > target {
			right = right - 1
		} else if arr[mid] < target {
			left = left + 1
		} else if arr[mid] == target {
			right = right - 1
		}
		if left > right {
			break
		}
	}
	if left >= len(arr) || arr[left] != target {
		return -1
	}
	return left
}
