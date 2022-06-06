package code0773

import (
	"fmt"
	"strings"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	board := [][]int{{4, 1, 2}, {5, 0, 3}}
	target := "123450"
	r := exection(board, target)
	fmt.Println(r)
}

func exection(board [][]int, target string) int {
	neighbor := [][]int{
		{1, 3},
		{0, 4, 2},
		{1, 5},
		{0, 4},
		{3, 1, 5},
		{4, 2},
	}
	m := 2
	n := 3
	var sbuilder strings.Builder
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			sbuilder.WriteString(fmt.Sprint(board[i][j]))
		}
	}
	start := sbuilder.String()

	visited := common.NewHashSet()
	queue := common.NewQueue()
	step := 0
	visited.Add(start)
	queue.Push(start)
	for {
		if queue.Empty() {
			break
		}
		n := queue.Len()
		for i := 0; i < n; i++ {
			curr := queue.Pop().(string)
			if curr == target {
				return step
			}
			index := 0
			for {
				if curr[index] == '0' {
					break
				}
				index++
			}
			for _, v := range neighbor[index] {
				now := change(curr, v, index)
				if !visited.Contains(now) {
					visited.Add(now)
					queue.Push(now)
				}
			}
		}
		step++
	}
	return -1
}

func change(curr string, v, index int) string {
	t := []byte(curr)
	t[v], t[index] = t[index], t[v]
	return string(t)
}
