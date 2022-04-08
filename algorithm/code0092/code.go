package code0092

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	list := common.NewLinkList()
	left := 1
	right := 3
	r := exection(&list, left, right)
	r.Display()
}

var node *common.LinkList

func exection(list *common.LinkList, left, right int) *common.LinkList {
	if left == 1 {
		return reverse(list, right)
	}
	list.Next = exection(list, left-1, right-1)
	return list
}

func reverse(list *common.LinkList, n int) *common.LinkList {
	if n == 1 {
		node = list.Next
		return list
	}
	last := reverse(list.Next, n-1)
	list.Next.Next = list
	list.Next = node
	return last
}
