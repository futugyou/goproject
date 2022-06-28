package code0174

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	grid := [][]int{{-2, -3, 3}, {-5, -10, 1}, {10, 30, -5}}
	r := exection(grid)
	fmt.Println(r)
}

var memo [][]int

func exection(grid [][]int) int {
	m := len(grid)
	n := len(grid[0])
	memo = make([][]int, m)
	for i := 0; i < m; i++ {
		memo[i] = make([]int, n)
		for j := 0; j < n; j++ {
			memo[i][j] = -1
		}
	}
	return dp(grid, 0, 0)
}

func dp(grid [][]int, i, j int) int {
	m := len(grid)
	n := len(grid[0])
	if i == m-1 && j == n-1 {
		if grid[i][j] >= 0 {
			return 1
		} else {
			return 1 - grid[i][j]
		}

	}
	if i == m || j == n {
		return math.MaxInt
	}
	if memo[i][j] != -1 {
		return memo[i][j]
	}
	res := common.Min(dp(grid, i, j+1), dp(grid, i+1, j)) - grid[i][j]
	if res <= 0 {
		memo[i][j] = 1
	} else {
		memo[i][j] = res
	}
	return memo[i][j]
}
