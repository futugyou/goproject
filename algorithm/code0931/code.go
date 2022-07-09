package code0931

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	matrix := [][]int{{2, 1, 3}, {6, 5, 4}, {7, 8, 9}}
	r := exection(matrix)
	fmt.Println(r)
}

var memo [][]int

func exection(matrix [][]int) int {
	n := len(matrix)
	memo = make([][]int, n)
	for i := 0; i < n; i++ {
		memo[i] = make([]int, n)
		for j := 0; j < n; j++ {
			memo[i][j] = 100
		}
	}
	res := 100
	for i := 0; i < n; i++ {
		res = common.Min(res, dp(matrix, n-1, i))
	}
	return res
}

func dp(matrix [][]int, i, j int) int {
	n := len(matrix)
	if i < 0 || j < 0 || i >= n || j >= n {
		return 100
	}
	if i == 0 {
		return matrix[0][j]
	}
	if memo[i][j] != 100 {
		return memo[i][j]
	}
	memo[i][j] = matrix[i][j] + common.Min(dp(matrix, i-1, j), common.Min(dp(matrix, i-1, j-1), dp(matrix, i-1, j+1)))
	return memo[i][j]
}
