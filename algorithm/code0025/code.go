package code0025

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	list := common.NewLinkList()
	k := 3
	r := exection(&list, k)
	r.Display()
}

func exection(head *common.LinkList, k int) *common.LinkList {
	if head == nil {
		return nil
	}
	a := head
	b := head
	for i := 0; i < k; i++ {
		if b == nil {
			return head
		}
		b = b.Next
	}
	newHead := reverse(a, b)
	a.Next = exection(b, k)
	return newHead
}

func reverse(a, b *common.LinkList) *common.LinkList {
	var pre *common.LinkList = nil
	curr := a
	next := a
	for {
		if curr == b {
			break
		}
		next = curr.Next
		curr.Next = pre
		pre = curr
		curr = next
	}
	return pre
}
