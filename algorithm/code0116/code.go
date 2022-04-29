package code0116

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	root := common.NewBSTTreeNode()
	r := exection(&root)
	r.Display()
}

func exection(root *common.TreeNode) *common.TreeNode {
	if root == nil {
		return nil
	}
	Link(root.Left, root.Right)
	return root
}

func Link(left *common.TreeNode, right *common.TreeNode) {
	if left == nil || right == nil {
		return
	}
	left.Next = right
	Link(left.Left, left.Right)
	Link(right.Left, right.Right)
	Link(left.Right, right.Left)
}
