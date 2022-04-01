package code0235

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	node := common.NewBSTTreeNode()
	a := 3
	b := 6
	r := find(&node, a, b)
	if r != nil {
		fmt.Println(r.Val)
	}
}

func find(node *common.TreeNode, a, b int) *common.TreeNode {
	if node == nil {
		return nil
	}
	if node.Val < a {
		return find(node.Left, a, b)
	}
	if node.Val > b {
		return find(node.Right, a, b)
	}
	return node
}
