package code0141

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
	fmt.Println(r)
}
