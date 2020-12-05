package main

import (
	"container/heap"
	"fmt"
)

type heapFloat32 []float32

func (n *heapFloat32) Pop() interface{} {
	old := *n
	x := old[len(old)-1]
	new := old[0 : len(old)-1]
	*n = new
	return x
}

func (n *heapFloat32) Push(x interface{}) {
	*n = append(*n, x.(float32))
}

func (n heapFloat32) Len() int {
	return len(n)
}

func (n heapFloat32) Less(a, b int) bool {
	return n[a] < n[b]
}

func (n heapFloat32) Swap(a, b int) {
	n[a], n[b] = n[b], n[a]
}

func main() {
	myheap := &heapFloat32{2, 4.2, 1, 1.2, 3}
	heap.Init(myheap)
	size := len(*myheap)
	fmt.Println("size ", size)
	fmt.Printf("%v\n", myheap)
	myheap.Push(float32(-9))
	myheap.Push(float32(-1))
	myheap.Push(float32(-4))
	fmt.Printf("%v\n", myheap)
	heap.Init(myheap)
	fmt.Printf("%v\n", myheap)
}
