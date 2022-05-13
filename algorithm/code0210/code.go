package code0210

import "fmt"

func Exection() {
	prerequisites := [][]int{{1, 0}}
	numCourses := 2
	exection(prerequisites, numCourses)
}

func exection(dislikes [][]int, n int) {
	onPath = make([]bool, n)
	visited = make([]bool, n)
	postorder = make([]int, 0)
	graph := buildGraph(dislikes, n)
	for i := 0; i < n; i++ {
		build(graph, i)
	}
	if hasCycle {
		return
	}
	// 将后序遍历的结果进行反转（逆后序遍历顺序），就是拓扑排序的结果
	for i, j := 0, len(postorder)-1; i < j; i, j = i+1, j-1 {
		postorder[i], postorder[j] = postorder[j], postorder[i]
	}
	fmt.Println(postorder)
}

func buildGraph(dislikes [][]int, n int) [][]int {
	graph := make([][]int, n)
	for _, edge := range dislikes {
		a := edge[0]
		b := edge[1]
		graph[b] = append(graph[b], a)
	}
	fmt.Println(graph)
	return graph
}

func build(graph [][]int, i int) {
	if onPath[i] {
		hasCycle = true
	}
	if hasCycle || visited[i] {
		return
	}
	visited[i] = true
	onPath[i] = true
	for _, v := range graph[i] {
		build(graph, v)
	}
	postorder = append(postorder, i)
	onPath[i] = false
}

var postorder []int
var hasCycle bool
var onPath []bool
var visited []bool
