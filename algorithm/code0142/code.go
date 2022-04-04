package code0142

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	list := common.NewLinkList()
	exection(list)
}

func exection(list common.LinkList) {
	fast := &list
	slow := &list
	r := false
	for {
		if fast == nil || fast.Next == nil {
			break
		}
		fast = fast.Next.Next
		slow = slow.Next
		if fast == slow {
			r = true
			break
		}
	}
	if !r {
		fmt.Println("nil")
		return
	}
	fast = &list
	for {
		if fast == slow {
			fmt.Println(fast.Val)
			return
		}
		fast = fast.Next
		slow = slow.Next
	}

}
