package code0886

import "fmt"

func Exection() {
	dislikes := [][]int{{1, 2}, {1, 3}, {2, 4}}
	n := 4
	exection(dislikes, n)
}

func exection(dislikes [][]int, n int) {

	color = make([]bool, n+1)
	visited = make([]bool, n+1)
	graph := buildGraph(dislikes, n)
	for i := 1; i <= n; i++ {
		if !visited[i] {
			build(graph, i)
		}
	}
	fmt.Println(result)
}

func buildGraph(dislikes [][]int, n int) [][]int {
	graph := make([][]int, n+1)
	for _, edge := range dislikes {
		a := edge[0]
		b := edge[1]
		graph[a] = append(graph[a], b)
		graph[b] = append(graph[b], a)
	}
	fmt.Println(graph)
	return graph
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
