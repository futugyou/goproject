package code0743

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	times := [][]int{{2, 1, 1}, {2, 3, 1}, {3, 4, 1}}
	n := 4
	k := 2
	r := exection(times, n, k)
	fmt.Println(r)
}

func exection(times [][]int, n, k int) int {
	graph := make([][][]int, n+1)
	for i := 1; i <= n; i++ {
		graph[i] = make([][]int, 0)
	}

	for _, edge := range times {
		from := edge[0]
		to := edge[1]
		weight := edge[2]
		graph[from] = append(graph[from], []int{to, weight})
	}
	fmt.Println(graph)
	dijkstra := common.NewDijkstra(k, graph)
	distTo := dijkstra.ExecDist()
	fmt.Println(distTo)
	result := 0
	for i := 1; i < len(distTo); i++ {
		if distTo[i] == math.MaxInt {
			return -1
		}
		result = max(result, distTo[i])
	}
	return result
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
