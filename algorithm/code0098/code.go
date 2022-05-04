package code0098

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	node := common.NewTreeNode()
	r := exection(&node)
	fmt.Println(r)
	node = common.NewBSTTreeNode()
	r = exection(&node)
	fmt.Println(r)
}

func exection(root *common.TreeNode) bool {
	return check(root, nil, nil)
}

func check(root, min, max *common.TreeNode) bool {
	if root == nil {
		return true
	}
	if min != nil && min.Val >= root.Val {
		return false
	}
	if max != nil && max.Val <= root.Val {
		return false
	}
	return check(root.Left, min, root) && check(root.Right, root, max)
}
