package code0054

import (
	"fmt"
)

func Exection() {
	matrix := [][]int{{2, 3, 4},
		{5, 6, 7},
		{8, 9, 10},
		{11, 12, 13},
		{14, 15, 16}}
	r := exection(matrix)
	fmt.Println(r)
}

func exection(matrix [][]int) []int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return []int{}
	}
	if len(matrix) == 1 {
		return matrix[0]
	}
	return order(matrix, 0, len(matrix[0])-1, 0, len(matrix)-1, []int{})
}

func order(matrix [][]int, left, right, up, down int, nums []int) []int {
	if left > right || up > down {
		return nums
	}

	for i := left; i <= right; i++ {
		nums = append(nums, matrix[up][i])
	}

	stop := true
	for i := up + 1; i <= down; i++ {
		stop = false
		nums = append(nums, matrix[i][right])
	}
	if stop {
		return nums
	}

	stop = true
	for i := right - 1; i >= left; i-- {
		stop = false
		nums = append(nums, matrix[down][i])
	}
	if stop {
		return nums
	}

	for i := down - 1; i > up; i-- {
		if right > 0 {
			nums = append(nums, matrix[i][left])
		}
	}

	return order(matrix, left+1, right-1, up+1, down-1, nums)
}
