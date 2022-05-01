package code0222

import (
	"fmt"
	"math"

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
	left := root.Left
	leftcount := 0
	right := root.Right
	rightcount := 0
	for {
		if left == nil {
			break
		}
		left = left.Left
		leftcount++
	}
	for {
		if right == nil {
			break
		}
		right = right.Right
		rightcount++
	}
	if leftcount == rightcount {
		return int(math.Pow(2, float64(leftcount))) - 1
	}
	return 1 + exection(root.Left) + exection(root.Right)
}
