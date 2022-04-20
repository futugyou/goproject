package common

import (
	"container/heap"
	"fmt"
)

// Item 是优先队列中包含的元素。
type PriorityQueueItem struct {
	Value    int // 元素的值
	Priority int // 元素在队列中的优先级。
	// 元素的索引可以用于更新操作，它由 heap.Interface 定义的方法维护。
	index int // 元素在堆中的索引。
}

// 一个实现了 heap.Interface 接口的优先队列，队列中包含任意多个 Item 结构。
type PriorityQueue struct {
	items []*PriorityQueueItem
	Big   bool
}

func NewPriorityQueue(isBig bool) *PriorityQueue {
	pq := PriorityQueue{Big: isBig}
	heap.Init(&pq)
	return &pq
}
func (pq PriorityQueue) Len() int { return len(pq.items) }

func (pq PriorityQueue) Less(i, j int) bool {
	if pq.Big {
		return pq.items[i].Priority > pq.items[j].Priority
	} else {
		return pq.items[i].Priority < pq.items[j].Priority
	}
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(pq.items)
	item := x.(*PriorityQueueItem)
	item.index = n
	pq.items = append(pq.items, item)
	//heap.Push(pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old.items)
	item := old.items[n-1]
	item.index = -1
	pq.items = old.items[0 : n-1]
	//_ = heap.Pop(pq).(*PriorityQueue)
	return item
}

func (pq *PriorityQueue) Peek() interface{} {
	n := len(pq.items)
	item := pq.items[n-1]
	return item
}

func (pq *PriorityQueue) Size() int {
	n := len(pq.items)
	return n
}

// 更新函数会修改队列中指定元素的优先级以及值。
func (pq *PriorityQueue) Update(item *PriorityQueueItem, value int, priority int) {
	item.Value = value
	item.Priority = priority
	heap.Fix(pq, item.index)
}

// 这个示例首先会创建一个优先队列，并在队列中包含一些元素
// 接着将一个新元素添加到队列里面，并对其进行操作
// 最后按优先级有序地移除队列中的各个元素。
func PriorityQueueTest() {
	// 一些元素以及它们的优先级。
	items := map[int]int{
		3: 3, 2: 2, 4: 4,
	}

	// 创建一个优先队列，并将上述元素放入到队列里面，
	// 然后对队列进行初始化以满足优先队列（堆）的不变性。
	pq := PriorityQueue{
		Big:   false,
		items: make([]*PriorityQueueItem, len(items)),
	}
	i := 0
	for value, priority := range items {
		pq.items[i] = &PriorityQueueItem{
			Value:    value,
			Priority: priority,
			index:    i,
		}
		i++
	}
	heap.Init(&pq)

	// 插入一个新元素，然后修改它的优先级。
	item := &PriorityQueueItem{
		Value:    1,
		Priority: 1,
	}
	heap.Push(&pq, item)
	pq.Update(item, item.Value, 5)

	// 以降序形式取出并打印队列中的所有元素。
	for pq.Len() > 0 {
		item := heap.Pop(&pq).(*PriorityQueueItem)
		fmt.Printf("%d:%d ", item.Priority, item.Value)
	}
	fmt.Println()
}
