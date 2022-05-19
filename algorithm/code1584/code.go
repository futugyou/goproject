package code1584

import (
	"fmt"
	"math"
	"sort"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	points := [][]int{{0, 0}, {2, 2}, {3, 10}, {5, 2}, {7, 0}}
	r := exection(points)
	fmt.Println(r)
}

func exection(points [][]int) int {
	n := len(points)

	edges := make([][]int, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			xi := (float64)(points[i][0])
			yi := (float64)(points[i][1])
			xj := (float64)(points[j][0])
			yj := (float64)(points[j][1])
			edges = append(edges, []int{i, j, int(math.Abs(xi-xj)) + int(math.Abs(yi-yj))})
		}
	}
	fmt.Println(edges)
	sort.Slice(edges, func(i, j int) bool {
		return edges[i][2] < edges[j][2]
	})
	fmt.Println(edges)
	result := 0
	uf := common.NewUnionFind(n)
	for _, v := range edges {
		a := v[0]
		b := v[1]
		if uf.Connected(a, b) {
			continue
		}
		result += v[2]
		uf.Union(a, b)
	}
	return result
}
