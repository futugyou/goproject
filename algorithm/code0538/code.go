package code0538

import "github.com/futugyousuzu/goproject/algorithm/common"

var sum = 0

func Exection() {
	root := common.NewBSTTreeNode()
	r := exection(&root)
	r.Display()
}

func exection(root *common.TreeNode) *common.TreeNode {
	if root == nil {
		return nil
	}
	// left -> right asc
	// right -> left desc
	exection(root.Right)
	sum += root.Val
	root.Val = sum
	exection(root.Left)
	return root
}
