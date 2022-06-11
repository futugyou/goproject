package code0329

import (
	"fmt"
	"sort"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := [][]int{{3, 4, 5}, {3, 2, 6}, {2, 2, 1}}
	exection(nums)
}

var dir = [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

type item struct {
	i int
	j int
	v int
}

func exection(matrix [][]int) {
	m := len(matrix)
	n := len(matrix[0])
	dp := make([][]int, m)
	data := make([]item, 0)
	for i := 0; i < m; i++ {
		dp[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dp[i][j] = 1
			data = append(data, item{i: i, j: j, v: matrix[i][j]})
		}
	}
	sort.Slice(data, func(i, j int) bool {
		return data[i].v > data[j].v
	})
	for _, d := range data {
		for _, v := range dir {
			x, y := v[0]+d.i, v[1]+d.j
			if x >= 0 && x < m && y >= 0 && y < n {
				if matrix[d.i][d.j] < matrix[x][y] {
					dp[d.i][d.j] = common.Max(dp[d.i][d.j], dp[x][y]+1)
				}
			}
		}
	}
	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			res = common.Max(res, dp[i][j])
		}
	}
	fmt.Println(res)
}
