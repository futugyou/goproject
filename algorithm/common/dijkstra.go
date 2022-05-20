package common

import "math"

type Dijkstra struct {
	start int
	graph [][][]int
}

func NewDijkstra(start int, graph [][][]int) *Dijkstra {
	d := Dijkstra{}
	d.start = start
	d.graph = graph
	return &d
}

func (d *Dijkstra) ExecDist() []int {
	// 定义：distTo[i] 的值就是起点 start 到达节点 i 的最短路径权重
	distTo := make([]int, len(d.graph))
	for i := 0; i < len(distTo); i++ {
		distTo[i] = math.MaxInt
	}
	// base case，start 到 start 的最短距离就是 0
	distTo[d.start] = 0

	// 优先级队列，distFromStart 较小的排在前面
	pq := NewPriorityQueue2(stateCmpASC)
	// 从起点 start 开始进行 BFS
	pq.Push(newState(d.start, 0))

	for {
		if pq.Empty() {
			break
		}
		curState := pq.Pop().(*state)
		curNodeID := curState.id
		curDistFromStart := curState.distFromStart

		if curDistFromStart > distTo[curNodeID] {
			continue
		}
		// 将 curNode 的相邻节点装入队列
		for _, neighbor := range d.graph[curNodeID] {
			nextNodeID := neighbor[0]
			distToNextNode := distTo[curNodeID] + neighbor[1]
			// 更新 dp table
			if distTo[nextNodeID] > distToNextNode {
				distTo[nextNodeID] = distToNextNode
				pq.Push(newState(nextNodeID, distToNextNode))
			}
		}
	}
	return distTo
}

type state struct {
	// 图节点的 id
	id int
	// 从 start 节点到当前节点的距离
	distFromStart int
}

func newState(id, distFromStart int) *state {
	return &state{id: id, distFromStart: distFromStart}
}

func stateCmpASC(a, b interface{}) int {
	if a == b {
		return 0
	}
	x := a.(*state)
	y := b.(*state)
	if x.distFromStart > y.distFromStart {
		return 1
	} else if x.distFromStart < y.distFromStart {
		return -1
	}
	return 0
}
