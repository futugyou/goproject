package code0226

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

	t := root.Left
	root.Left = root.Right
	root.Right = t

	exection(root.Left)
	exection(root.Right)

	return root
}
