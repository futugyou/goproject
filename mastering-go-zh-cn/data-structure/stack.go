package main

import (
	"fmt"
)

type Node struct {
	Value int
	Next  *Node
}

var size = 0
var root = new(Node)

func push(v int) bool {
	if root == nil {
		root = &Node{v, nil}
		size = 1
		return true
	}
	tmp := &Node{v, nil}
	tmp.Next = root
	root = tmp
	size++
	return true
}

func pop(t *Node) (int, bool) {
	if size == 0 {
		return 0, false
	}
	if size == 1 {
		size--
		root = nil
		return t.Value, true
	}
	root = root.Next
	size--
	return t.Value, true
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
	root = nil
	pop(root)
	for i := 0; i < 10; i++ {
		push(i)
	}
	traverse(root)
	pop(root)
	traverse(root)
}
