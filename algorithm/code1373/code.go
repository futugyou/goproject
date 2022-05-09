package code1373

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	root := common.NewBSTTreeNode()
	exection(&root)
	fmt.Println(maxsum)
}

var maxsum = 0

func exection(root *common.TreeNode) []int {
	if root == nil {
		// [0]: 1 bst, 0 not bst
		// [1]: min val in root's nodes
		// [2]: max val in root's nodes
		// [3]: sum in root's nodes
		return []int{1, math.MaxInt, math.MinInt, 0}
	}
	left := exection(root.Left)
	right := exection(root.Right)
	result := make([]int, 4)
	if left[0] == 1 && right[0] == 1 && left[2] < root.Val && root.Val < right[1] {
		result[0] = 1
		result[1] = min(left[1], root.Val)
		result[2] = max(right[2], root.Val)
		result[3] = left[3] + right[3] + root.Val
		maxsum = max(maxsum, result[3])
	}
	return result
}

func min(i1, i2 int) int {
	if i1 < i2 {
		return i1
	}
	return i2
}

func max(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}
