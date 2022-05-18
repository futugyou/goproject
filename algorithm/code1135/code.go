package code1135

import (
	"fmt"
	"sort"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	conections := [][]int{{1, 2, 5}, {1, 3, 6}, {2, 3, 1}}
	n := 3
	r := exection(conections, n)
	fmt.Println(r)
}

func exection(conections [][]int, n int) int {
	uf := common.NewUnionFind(n + 1)
	fmt.Println(conections)
	sort.Slice(conections, func(i, j int) bool {
		return conections[i][2] < conections[j][2]
	})
	fmt.Println(conections)
	result := 0
	for _, v := range conections {
		a := v[0]
		b := v[1]
		if uf.Connected(a, b) {
			continue
		}
		result += v[2]
		uf.Union(a, b)
	}
	if uf.Count() == 2 {
		return result
	}
	return -1
}
