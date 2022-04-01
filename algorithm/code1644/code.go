package code1644

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	node := common.NewTreeNode()
	a := 3
	b := 6
	r := find(&node, a, b)
	if !finda || !findb || r == nil {
		fmt.Println("nil")
		return
	}

	fmt.Println(r.Val)
}

var (
	finda bool = false
	findb bool = false
)

func find(node *common.TreeNode, a, b int) *common.TreeNode {
	if node == nil {
		return nil
	}
	left := find(node.Left, a, b)
	right := find(node.Right, a, b)
	if left != nil && right != nil {
		return node
	}

	if node.Val == a {
		finda = true
		return node
	}

	if node.Val == b {
		findb = true
		return node
	}

	if left == nil {
		return right
	}
	return left
}
