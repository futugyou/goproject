package code0695

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	grid := [][]int{}
	exection(grid)
}

func exection(grid [][]int) {
	m := len(grid)
	n := len(grid[0])
	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				res = common.Max(res, dfs(grid, i, j))
			}
		}
	}
	fmt.Println(res)
}

func dfs(grid [][]int, i, j int) int {
	m := len(grid)
	n := len(grid[0])
	if i < 0 || j < 0 || i >= m || j >= n {
		return 0
	}
	if grid[i][j] == 0 {
		return 0
	}
	grid[i][j] = 0

	return dfs(grid, i+1, j) +
		dfs(grid, i-1, j) +
		dfs(grid, i, j+1) +
		dfs(grid, i, j-1) + 1

}
