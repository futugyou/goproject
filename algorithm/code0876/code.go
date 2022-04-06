package code0876

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	list := common.NewLinkList()
	exection(list)
}

func exection(list common.LinkList) {
	fast := &list
	slow := &list
	for {
		if fast == nil || fast.Next == nil {
			break
		}
		fast = fast.Next.Next
		slow = slow.Next
	}
	slow.Display()
}
