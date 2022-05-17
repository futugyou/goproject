package code0261

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	edges := [][]int{{0, 1}, {1, 2}, {2, 3}, {1, 3}, {1, 4}}
	n := 5
	r := exection(edges, n)
	fmt.Println(r)
}

func exection(edges [][]int, n int) bool {
	uf := common.NewUnionFind(n)
	for _, v := range edges {
		a := v[0]
		b := v[1]
		if uf.Connected(a, b) {
			return false
		}
		uf.Union(a, b)
	}
	return uf.Count() == 1
}
