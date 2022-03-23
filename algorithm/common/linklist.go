package common

type LinkList struct {
	Val  int
	Next *LinkList
}

func NewLinkList() LinkList {
	h := LinkList{
		Val: 1,
		Next: &LinkList{
			Val: 1,
			Next: &LinkList{
				Val: 2,
				Next: &LinkList{
					Val: 2,
					Next: &LinkList{
						Val: 3,
						Next: &LinkList{
							Val: 4,
						},
					},
				},
			},
		},
	}
	return h
}
