package code0077

import "fmt"

func Exection() {
	n := 4
	k := 2
	exection(n, k)
	fmt.Println(result)
}

var result [][]int

func exection(n, k int) {
	result = make([][]int, 0)
	path := make([]int, 0)
	backtrack(n, k, 1, path)
}

func backtrack(n, k, start int, path []int) {
	if len(path) == k {
		t := make([]int, 2)
		copy(t, path)
		result = append(result, t)
		fmt.Println(result)
		return
	}
	for i := start; i <= n; i++ {
		path = append(path, i)
		backtrack(n, k, i+1, path)
		path = path[:len(path)-1]
	}
}
