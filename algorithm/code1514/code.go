package code1514

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	edges := [][]int{{0, 1}, {1, 2}, {0, 2}}
	succProb := []float32{0.5, 0.5, 0.2}
	n := 3
	start := 0
	end := 2
	r := exection(edges, succProb, n, start, end)
	fmt.Println(r)
}

type graphItem struct {
	to     int
	weight float32
}

type state struct {
	id            int
	probFromStart float32
}

func exection(edges [][]int, succProb []float32, n, start, end int) float32 {
	graph := make([][]graphItem, n)
	for i := 0; i < len(edges); i++ {
		from := edges[i][0]
		to := edges[i][1]
		weight := succProb[i]
		graph[from] = append(graph[from], graphItem{to, weight})
	}
	probTo := make([]float32, n)
	for i := 0; i < n; i++ {
		probTo[i] = -1
	}
	probTo[start] = 1
	pq := common.NewPriorityQueue2(stateCmpDESC)
	pq.Push(state{id: start, probFromStart: 1})
	for {
		if pq.Empty() {
			break
		}
		curr := pq.Pop().(state)
		id := curr.id
		probFromStart := curr.probFromStart
		if id == end {
			return probFromStart
		}
		if probFromStart < probTo[id] {
			continue
		}
		for _, neighbor := range graph[id] {
			nextId := neighbor.to
			probToNextNode := probTo[id] * neighbor.weight
			if probToNextNode > probTo[nextId] {
				probTo[nextId] = probToNextNode
				pq.Push(state{id: nextId, probFromStart: probToNextNode})
			}
		}
	}
	return 0
}

func stateCmpDESC(a, b interface{}) int {
	if a == b {
		return 0
	}
	x := a.(state)
	y := b.(state)
	if x.probFromStart > y.probFromStart {
		return -1
	} else if x.probFromStart < y.probFromStart {
		return 1
	}
	return 0
}
