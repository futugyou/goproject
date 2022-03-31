package common

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func NewTreeNode() TreeNode {
	h := TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val: 2,
			Left: &TreeNode{
				Val: 3,
			},
		},
		Right: &TreeNode{
			Val: 4,
			Right: &TreeNode{
				Val: 5,
				Left: &TreeNode{
					Val: 6,
				},
			},
		},
	}
	return h
}
