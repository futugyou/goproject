package code0111

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	root := common.NewBSTTreeNode()
	r := exection(&root)
	fmt.Println(r)
}

func exection(root *common.TreeNode) int {
	if root == nil {
		return 0
	}
	result := 1
	queue := common.NewQueue()
	queue.Push(root)
	for {
		n := queue.Len()
		if n == 0 {
			break
		}
		for i := 0; i < n; i++ {
			curr := queue.Pop().(*common.TreeNode)
			if curr.Left == nil && curr.Right == nil {
				return result
			}
			if curr.Left != nil {
				queue.Push(curr.Left)
			}
			if curr.Right != nil {
				queue.Push(curr.Right)
			}
		}
		result++
	}
	return result
}
