package code0059

import (
	"fmt"
)

func Exection() {
	n := 4
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) [][]int {
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
	}
	return order(matrix, 0, n-1, 0, n-1, 1)
}

func order(matrix [][]int, left, right, up, down int, index int) [][]int {
	if left > right || up > down {
		return matrix
	}
	for i := left; i <= right; i++ {
		matrix[up][i] = index
		index++
	}

	stop := true
	for i := up + 1; i <= down; i++ {
		stop = false
		matrix[i][right] = index
		index++
	}
	if stop {
		return matrix
	}

	stop = true
	for i := right - 1; i >= left; i-- {
		stop = false
		matrix[down][i] = index
		index++
	}
	if stop {
		return matrix
	}

	for i := down - 1; i > up; i-- {
		if right > 0 {
			matrix[i][left] = index
			index++
		}
	}

	return order(matrix, left+1, right-1, up+1, down-1, index)
}
