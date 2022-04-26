package code0654

import (
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{3, 2, 1, 6, 0, 5}
	r := exection(nums)
	r.Display()
}

func exection(nums []int) *common.TreeNode {
	return build(nums, 0, len(nums)-1)
}

func build(nums []int, lo, hi int) *common.TreeNode {
	if lo > hi {
		return nil
	}
	index := -1
	maxval := math.MinInt
	for i := lo; i <= hi; i++ {
		if maxval < nums[i] {
			index = i
			maxval = nums[i]
		}
	}

	root := &common.TreeNode{Val: maxval}
	root.Left = build(nums, lo, index-1)
	root.Right = build(nums, index+1, hi)
	return root
}
