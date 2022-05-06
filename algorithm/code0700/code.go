package code0700

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	node := common.NewBSTTreeNode()
	key := 1
	r := exection(&node, key)
	if r == nil {
		fmt.Println("null")
	} else {
		fmt.Println(r.Val)
	}
}

func exection(root *common.TreeNode, key int) *common.TreeNode {
	if root == nil {
		return nil
	}
	if root.Val < key {
		return exection(root.Right, key)
	}

	if root.Val > key {
		return exection(root.Left, key)
	}
	return root
}
