package code1631

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	heights := [][]int{{1, 2, 2}, {3, 8, 2}, {5, 3, 5}}

	r := exection(heights)
	fmt.Println(r)
}

type state struct {
	x             int
	y             int
	probFromStart int
}

func exection(heights [][]int) int {
	m := len(heights)
	n := len(heights[0])

	effortTo := make([][]int, n)
	for i := 0; i < m; i++ {
		effortTo[i] = make([]int, n)
		for j := 0; j < n; j++ {
			effortTo[i][j] = math.MaxInt
		}
	}
	effortTo[0][0] = 1
	pq := common.NewPriorityQueue2(stateCmpASC)

	pq.Push(state{x: 0, y: 0, probFromStart: 1})
	for {
		if pq.Empty() {
			break
		}
		curr := pq.Pop().(state)
		x := curr.x
		y := curr.y
		probFromStart := curr.probFromStart
		if x == m-1 && y == n-1 {
			return probFromStart
		}
		if probFromStart > effortTo[x][y] {
			continue
		}
		for _, neighbor := range getNeighbor(heights, x, y) {
			xx := neighbor[0]
			yy := neighbor[1]
			probToNextNode := max(effortTo[x][y], (int)(math.Abs(float64(heights[x][y]-heights[xx][yy]))))
			if probToNextNode < effortTo[xx][yy] {
				effortTo[xx][yy] = probToNextNode
				pq.Push(state{x: xx, y: yy, probFromStart: probToNextNode})
			}
		}
	}
	return 0
}

func max(i int, f int) int {
	if i > f {
		return i
	}
	return f
}

func getNeighbor(heights [][]int, x, y int) [][]int {
	m := len(heights)
	n := len(heights[0])
	result := make([][]int, 0)
	for _, dir := range dirs {
		xx := x + dir[0]
		yy := y + dir[1]
		if xx < 0 || xx >= m || yy < 0 || yy >= n {
			continue
		}
		result = append(result, []int{xx, yy})
	}
	return result
}

var dirs = [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

func stateCmpASC(a, b interface{}) int {
	if a == b {
		return 0
	}
	x := a.(state)
	y := b.(state)
	if x.probFromStart < y.probFromStart {
		return -1
	} else if x.probFromStart > y.probFromStart {
		return 1
	}
	return 0
}
