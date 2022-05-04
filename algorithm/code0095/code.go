package code0095

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	n := 3
	r := exection(n)
	fmt.Println(len(r))
}

func exection(n int) []*common.TreeNode {
	if n <= 0 {
		return make([]*common.TreeNode, 0)
	}
	return build(1, n)
}

func build(lo, hi int) []*common.TreeNode {
	res := make([]*common.TreeNode, 0)
	if lo > hi {
		res = append(res, nil)
		return res
	}
	for i := lo; i <= hi; i++ {
		leftlist := build(lo, i-1)
		rightlist := build(i+1, hi)
		for _, left := range leftlist {
			for _, right := range rightlist {
				root := common.TreeNode{Val: i}
				root.Left = left
				root.Right = right
				res = append(res, &root)
			}
		}
	}
	return res
}
