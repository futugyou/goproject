package main

import (
	"fmt"
)

type Node struct {
	Value    int
	Previous *Node
	Next     *Node
}

var root = new(Node)

func addNode(t *Node, v int) int {
	if root == nil {
		t = &Node{v, nil, nil}
		root = t
		return 0
	}
	if v == t.Value {
		fmt.Println("already exusts", v)
		return -1
	}
	if t.Next == nil {
		tmp := t
		t.Next = &Node{v, tmp, nil}
		return -2
	}
	return addNode(t.Next, v)
}

func traverse(t *Node) {
	if t == nil {
		fmt.Println("-> empty")
		return
	}
	for t != nil {
		fmt.Printf("%d -> ", t.Value)
		t = t.Next
	}
	fmt.Println()
}

func reverse(t *Node) {
	if t == nil {
		fmt.Println("-> empty")
		return
	}

	tmp := t
	for t != nil {
		tmp = t
		t = t.Next
	}

	for tmp.Previous != nil {
		fmt.Printf("%d -> ", tmp.Value)
		tmp = tmp.Previous
	}
	fmt.Printf("%d -> ", tmp.Value)
	fmt.Println()
}

func size(t *Node) int {
	if t == nil {
		fmt.Println("-> empty")
		return 0
	}
	n := 0
	for t != nil {
		n++
		t = t.Next
	}
	return n
}

func lookupNode(t *Node, v int) bool {
	if root == nil {
		return false
	}
	if v == t.Value {
		return true
	}
	if t.Next == nil {
		return false
	}
	return lookupNode(t.Next, v)
}

func main() {
	root = nil
	traverse(root)
	addNode(root, 1)
	addNode(root, 12)
	addNode(root, 13)
	addNode(root, 14)
	traverse(root)
	reverse(root)
	traverse(root)
}
