package code0297

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	root := common.NewBSTTreeNode()
	r := exection(&root)
	r.Display()
	fmt.Println()
}

func exection(root *common.TreeNode) *common.TreeNode {
	resultstring := serialize(root)
	fmt.Println(resultstring)
	r := deserialize(resultstring)
	return r
}

func deserialize(data string) *common.TreeNode {
	sp := strings.Split(data, ",")
	var build func() *common.TreeNode
	build = func() *common.TreeNode {
		if sp[0] == "null" {
			sp = sp[1:]
			return nil
		}
		val, _ := strconv.Atoi(sp[0])
		sp = sp[1:]
		// if err != nil {
		// 	return nil
		// }
		return &common.TreeNode{Val: val, Left: build(), Right: build()}
		// root := &common.TreeNode{Val: val}
		// root.Left = build()
		// root.Right = build()
		// return root
	}
	return build()
}

func serialize(root *common.TreeNode) string {
	sb := &strings.Builder{}
	var dfs func(*common.TreeNode)
	dfs = func(node *common.TreeNode) {
		if node == nil {
			sb.WriteString("null,")
			return
		}
		sb.WriteString(strconv.Itoa(node.Val))
		sb.WriteByte(',')
		dfs(node.Left)
		dfs(node.Right)
	}
	dfs(root)
	return sb.String()
}
