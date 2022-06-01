package code0694

import (
	"fmt"
	"strings"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	grid := [][]int{}
	exection(grid)
}

func exection(grid [][]int) {
	m := len(grid)
	n := len(grid[0])
	set := common.NewHashSet()
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == 1 {
				var sbuilder strings.Builder
				dfs(grid, i, j, sbuilder, "0")
				set.Add(sbuilder.String())
			}
		}
	}
	fmt.Println(set.Size())
}

func dfs(grid [][]int, i, j int, sb strings.Builder, dir string) {
	m := len(grid)
	n := len(grid[0])
	if i < 0 || j < 0 || i >= m || j >= n {
		return
	}
	if grid[i][j] == 0 {
		return
	}
	grid[i][j] = 0
	sb.WriteString(dir)
	sb.WriteString(",")
	dfs(grid, i+1, j, sb, "1")
	dfs(grid, i-1, j, sb, "2")
	dfs(grid, i, j+1, sb, "3")
	dfs(grid, i, j-1, sb, "4")
	sb.WriteString("-")
	sb.WriteString(dir)
	sb.WriteString(",")
}
