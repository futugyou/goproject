package main

import (
	"fmt"
)

type Node struct {
	Value int
	Next  *Node
}

var size = 0
var queue = new(Node)

func push(t *Node, v int) bool {
	if queue == nil {
		queue = &Node{v, nil}
		size++
		return true
	}
	t = &Node{v, nil}
	t.Next = queue
	queue = t
	size++
	return true
}

func pop(t *Node) (int, bool) {
	if size == 0 {
		return 0, false
	}
	if size == 1 {
		queue = nil
		size--
		return t.Value, true
	}
	tmp := t
	for t.Next != nil {
		tmp = t
		t = t.Next
	}
	v := tmp.Next.Value
	tmp.Next = nil
	size--
	return v, true
}

func traverse(t *Node) {
	if size == 0 {
		fmt.Println("empty")
		return
	}
	for t != nil {
		fmt.Printf("%d -> ", t.Value)
		t = t.Next
	}
	fmt.Println()
}

func main() {
	queue = nil
	for i := 0; i < 10; i++ {
		push(queue, i)
	}
	traverse(queue)
	pop(queue)
	traverse(queue)
}
