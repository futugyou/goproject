package code0797

import "fmt"

func Exection() {
	graph := [][]int{{1, 2}, {3}, {3}, {}}
	exection(graph)
}

func exection(graph [][]int) {
	result = make([][]int, 0)
	path := make([]int, 0)
	build(graph, 0, path)
	for _, v := range result {
		fmt.Println(v)
	}
}

func build(graph [][]int, i int, path []int) {
	path = append(path, i)
	n := len(graph)
	if n-1 == i {
		result = append(result, path)
		return
	}
	for _, v := range graph[i] {
		build(graph, v, path)
	}
}

var result [][]int
