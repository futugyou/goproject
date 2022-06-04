package code1905

import (
	"fmt"
)

func Exection() {
	grid := [][]int{}
	grid2 := [][]int{}
	exection(grid, grid2)
}

func exection(grid, grid2 [][]int) {
	m := len(grid)
	n := len(grid[0])
	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 0 && grid2[i][j] == 1 {
				dfs(grid2, i, j)
			}
		}
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				res++
				dfs(grid2, i, j)
			}
		}
	}
	fmt.Println(res)
}

func dfs(grid [][]int, i, j int) {
	m := len(grid)
	n := len(grid[0])
	if i < 0 || j < 0 || i >= m || j >= n {
		return
	}
	if grid[i][j] == 0 {
		return
	}
	grid[i][j] = 0

	dfs(grid, i+1, j)
	dfs(grid, i-1, j)
	dfs(grid, i, j+1)
	dfs(grid, i, j-1)

}
