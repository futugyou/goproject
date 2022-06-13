package code0337

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	root := common.NewBSTTreeNode()
	exection(&root)
}

var memo map[*common.TreeNode]int

func exection(root *common.TreeNode) {
	memo = make(map[*common.TreeNode]int)
	r := rob(root)
	fmt.Println(r)
}

func rob(root *common.TreeNode) int {
	if root == nil {
		return 0
	}
	if v, ok := memo[root]; ok {
		return v
	}
	do := root.Val
	if root.Left != nil {
		do += rob(root.Left.Left) + rob(root.Left.Right)
	}
	if root.Right != nil {
		do += rob(root.Right.Left) + rob(root.Right.Right)
	}
	notdo := rob(root.Left) + rob(root.Right)
	result := common.Max(do, notdo)
	memo[root] = result
	return result
}
