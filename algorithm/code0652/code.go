package code0652

import (
	"fmt"
	"strconv"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

var dic map[string]int
var list []*common.TreeNode

func Exection() {
	dic = make(map[string]int)
	list = make([]*common.TreeNode, 0)
	root := common.NewBSTTreeNode()
	exection(&root)
	fmt.Println(len(list))
}

func exection(root *common.TreeNode) string {
	if root == nil {
		return "#"
	}
	left := exection(root.Left)
	right := exection(root.Right)
	sub := left + "," + strconv.Itoa(root.Val) + "," + right
	f := dic[sub]
	if f == 1 {
		list = append(list, root)
	}
	dic[sub] = f + 1
	return sub
}
