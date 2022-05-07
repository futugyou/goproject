package code0230

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	root := common.NewBSTTreeNode()
	index := 7
	exection(&root, index)
	fmt.Println(result)
}

var result = 0
var step = 0

func exection(root *common.TreeNode, index int) {
	if root == nil {
		return
	}
	exection(root.Left, index)
	step++
	if step == index {
		result = root.Val
		return
	}
	exection(root.Right, index)
}
