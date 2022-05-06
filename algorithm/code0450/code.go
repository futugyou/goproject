package code0450

import "github.com/futugyousuzu/goproject/algorithm/common"

// 如果要递归地插入或者删除二叉树节点，递归函数一定要有返回值，且返回值要被正确的接收。
func Exection() {
	node := common.NewBSTTreeNode()
	key := 5
	r := exection(&node, key)
	r.Display()
}

func exection(root *common.TreeNode, key int) *common.TreeNode {
	if root == nil {
		return nil
	}
	if root.Val == key {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}
		t := getMax(root.Left)
		root.Val = t.Val
		root.Left = exection(root.Left, t.Val)
	} else if root.Val < key {
		root.Right = exection(root.Right, key)
	} else {
		root.Left = exection(root.Left, key)
	}
	return root
}

func getMax(root *common.TreeNode) *common.TreeNode {
	for {
		if root.Right == nil {
			return root
		}
		root = root.Right
	}
}
