package common

import "fmt"

type TreeNode struct {
	Val    int
	Left   *TreeNode
	Next   *TreeNode
	Right  *TreeNode
	Parent *TreeNode
}

func (l *TreeNode) Display() {
	if l == nil {
		return
	}
	curr := l
	if curr.Left != nil {
		curr.Left.Display()
	}
	if curr.Right != nil {
		curr.Right.Display()
	}
	fmt.Print(curr.Val, " ")
}

func NewTreeNode() TreeNode {
	h := TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 3,
			},
		},
		Right: &TreeNode{
			Val: 4,
			Right: &TreeNode{
				Val: 5,
				Left: &TreeNode{
					Val: 6,
				},
			},
		},
	}
	return h
}

func NewBSTTreeNode() TreeNode {
	h := TreeNode{
		Val: 3,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 1,
			},
		},
		Right: &TreeNode{
			Val: 4,
			Right: &TreeNode{
				Val: 6,
				Left: &TreeNode{
					Val: 5,
				},
			},
		},
	}
	return h
}
