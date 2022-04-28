package code0114

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
	exection(root.Left)
	exection(root.Right)
	right := root.Right
	left := root.Left

	root.Left = nil
	root.Right = left
	p := root
	for {
		if p.Right == nil {
			break
		}
		p = p.Right
	}
	p.Right = right
	return root
}
