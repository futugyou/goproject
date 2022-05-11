package code0785

import "fmt"

func Exection() {
	graph := [][]int{{1, 2, 3}, {0, 2}, {0, 1, 3}, {0, 2}}
	exection(graph)
}

func exection(graph [][]int) {
	n := len(graph)
	color = make([]bool, n)
	visited = make([]bool, n)
	for i := 0; i < n; i++ {
		if !visited[i] {
			build(graph, i)
		}
	}
	fmt.Println(result)
}

func build(graph [][]int, i int) {
	if !result {
		return
	}
	visited[i] = true
	for _, v := range graph[i] {
		if visited[v] {
			if color[v] == color[i] {
				result = false
			}
		} else {
			color[v] = !color[i]
			build(graph, v)
		}
	}
}

var result bool = true
var color []bool
var visited []bool
