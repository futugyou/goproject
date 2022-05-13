package code0207

import "fmt"

func Exection() {
	prerequisites := [][]int{{1, 0}, {0, 1}}
	numCourses := 2
	exection(prerequisites, numCourses)
}

func exection(dislikes [][]int, n int) {
	onPath = make([]bool, n)
	visited = make([]bool, n)
	graph := buildGraph(dislikes, n)
	for i := 0; i < n; i++ {
		build(graph, i)
	}
	fmt.Println(!hasCycle)
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
	onPath[i] = false
}

var hasCycle bool
var onPath []bool
var visited []bool
