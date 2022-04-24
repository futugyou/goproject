package code0105

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	preorder := []int{3, 9, 20, 15, 7}
	inorder := []int{9, 3, 15, 20, 7}
	r := exection(preorder, inorder)
	r.Display()
}

func exection(preorder []int, inorder []int) *common.TreeNode {
	return build(preorder, 0, len(preorder)-1, inorder, 0, len(inorder)-1)
}

func build(preorder []int, preStart, preEnd int, inorder []int, inStart, inEnd int) *common.TreeNode {
	if preStart > preEnd {
		return nil
	}
	rootVal := preorder[preStart]
	index := 0
	for i := inStart; i <= inEnd; i++ {
		if rootVal == inorder[i] {
			index = i
			break
		}
	}
	leftsize := index - inStart
	root := &common.TreeNode{Val: rootVal}
	root.Left = build(preorder, preStart+1, preStart+leftsize, inorder, inStart, index-1)
	root.Right = build(preorder, preStart+leftsize+1, preEnd, inorder, index+1, inEnd)
	return root
}
