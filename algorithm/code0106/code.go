package code0106

import "github.com/futugyousuzu/goproject/algorithm/common"

func Exection() {
	inorder := []int{9, 3, 15, 20, 7}
	postorder := []int{9, 15, 7, 20, 3}
	r := exection(inorder, postorder)
	r.Display()
}

func exection(inorder []int, postorder []int) *common.TreeNode {
	return build(inorder, 0, len(inorder)-1, postorder, 0, len(postorder)-1)
}

func build(inorder []int, inStart, inEnd int, postorder []int, postStart, postEnd int) *common.TreeNode {
	if inStart > inEnd {
		return nil
	}
	rootVal := postorder[postEnd]
	index := 0
	for i := inStart; i <= inEnd; i++ {
		if rootVal == inorder[i] {
			index = i
			break
		}
	}
	leftsize := index - inStart
	root := &common.TreeNode{Val: rootVal}
	root.Left = build(inorder, inStart, index-1, postorder, postStart, postStart+leftsize-1)
	root.Right = build(inorder, index+1, inEnd, postorder, postStart+leftsize, postEnd-1)
	return root
}
