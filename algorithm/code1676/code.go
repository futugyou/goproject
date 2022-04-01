package code1676

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	node := common.NewTreeNode()
	nodes := make([]common.TreeNode, 0)
	nodes = append(nodes, common.NewTreeNode())
	r := find(&node, nodes)

	fmt.Println(r.Val)
}

func find(root *common.TreeNode, nodes []common.TreeNode) *common.TreeNode {
	if root == nil {
		return nil
	}
	for _, n := range nodes {
		if root.Val == n.Val {
			return root
		}
	}

	left := find(root.Left, nodes)
	right := find(root.Right, nodes)

	if left != nil && right != nil {
		return root
	}

	if left != nil {
		return left
	}
	return right
}
