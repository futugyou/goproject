package code0064

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	grid := [][]int{{1, 3, 1}, {1, 5, 1}, {4, 2, 1}}
	r := exection(grid)
	fmt.Println(r)
}

var dir [][]int = [][]int{{0, -1}, {-1, 0}}

func exection(grid [][]int) int {
	m := len(grid)
	n := len(grid[0])
	dp := make([][]int, m)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dp[i][j] = math.MaxInt
		}
	}
	dp[0][0] = grid[0][0]
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			for _, d := range dir {
				x, y := i+d[0], j+d[1]
				if x >= 0 && y >= 0 && x < m && y < n {
					dp[i][j] = common.Min(dp[i][j], dp[x][y]+grid[i][j])
				}
			}
		}
	}
	fmt.Println(dp)
	return dp[m-1][n-1]
}
