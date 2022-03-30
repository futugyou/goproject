package code0021

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	list1 := common.NewLinkList()
	list2 := common.NewLinkList()
	exection(list1, list2)
}

func exection(list1, list2 common.LinkList) {
	dummy := common.LinkList{Val: -1}
	p := &dummy
	l1 := &list1
	l2 := &list2
	for {
		if l1 == nil || l2 == nil {
			break
		}
		if l1.Val <= l2.Val {
			p.Next = l1
			l1 = l1.Next
		} else {
			p.Next = l2
			l2 = l2.Next
		}
		p = p.Next
	}
	if l1 != nil {
		p.Next = l1
	}
	if l2 != nil {
		p.Next = l2
	}
	curr := dummy.Next
	curr.Display()
}
