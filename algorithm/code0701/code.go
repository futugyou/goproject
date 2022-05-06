package code0701

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

// 如果要递归地插入或者删除二叉树节点，递归函数一定要有返回值，且返回值要被正确的接收。
func Exection() {
	node := common.NewBSTTreeNode()
	key := 7
	r := exection(&node, key)
	r.Display()
}

func exection(root *common.TreeNode, key int) *common.TreeNode {
	if root == nil {
		return &common.TreeNode{Val: key}
	}
	if root.Val < key {
		root.Right = exection(root.Right, key)
	}

	if root.Val > key {
		root.Left = exection(root.Left, key)
	}
	return root
}
