package common

import "fmt"

type LinkList struct {
	Val  int
	Next *LinkList
}

func (l LinkList) Less(other interface{}) bool {
	return l.Val < other.(LinkList).Val
}

func (l *LinkList) Display() {
	curr := l
	for {
		if curr == nil {
			break
		}
		fmt.Print(curr.Val, "")
		curr = curr.Next
	}
	fmt.Println("")
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
