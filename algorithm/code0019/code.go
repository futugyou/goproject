package code0019

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	list := common.NewLinkList()
	k := 2
	exection(list, k)
}

func exection(list common.LinkList, k int) {
	dummy := common.LinkList{Val: -1}
	dummy.Next = &list
	fast := &dummy
	slow := &dummy
	for i := 0; i <= k; i++ {
		if fast != nil {
			fast = fast.Next
		} else {
			return
		}
	}
	for {
		if fast == nil {
			break
		}
		fast = fast.Next
		slow = slow.Next
	}
	slow.Next = slow.Next.Next
	curr := dummy.Next
	curr.Display()
}
