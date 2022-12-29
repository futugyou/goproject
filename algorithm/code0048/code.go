package code0048

import "fmt"

func Exection() {
	matrix := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	r := exection(matrix)
	fmt.Println(r)
}

func exection(matrix [][]int) [][]int {
	n := len(matrix)
	if n <= 1 {
		return matrix
	}
	// 右侧对角线折叠
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
		}
	}

	// 左右折叠
	for i := 0; i < n; i++ {
		for j := 0; j < n/2; j++ {
			matrix[i][j], matrix[i][n-1-j] = matrix[i][n-1-j], matrix[i][j]
		}
	}
	return matrix
}
