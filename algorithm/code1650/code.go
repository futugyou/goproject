package code1650

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	node1 := common.NewTreeNode()
	node2 := common.NewTreeNode()
	r := find(&node1, &node2)

	fmt.Println(r.Val)
}

func find(treeNode1, treeNode2 *common.TreeNode) *common.TreeNode {
	a := treeNode1
	b := treeNode2
	for {
		if a == b {
			break
		}
		if a == nil {
			a = treeNode2
		} else {
			a = a.Parent
		}
		if b == nil {
			b = treeNode1
		} else {
			b = b.Parent
		}
	}
	return a
}
