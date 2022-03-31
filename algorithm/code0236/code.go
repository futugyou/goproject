package code0236

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	node := common.NewTreeNode()
	a := 3
	b := 6
	r := lowestCommonAncestor(&node, a, b)
	if r != nil {
		fmt.Println(r.Val)
	}
}

func lowestCommonAncestor(node *common.TreeNode, a, b int) *common.TreeNode {
	if node == nil {
		return nil
	}
	if node.Val == a || node.Val == b {
		return node
	}
	left := lowestCommonAncestor(node.Left, a, b)
	right := lowestCommonAncestor(node.Right, a, b)
	if left != nil && right != nil {
		return node
	}
	if left == nil {
		return right
	}
	return left
}
