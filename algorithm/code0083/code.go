package code0083

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	head := common.NewLinkList()
	exection(&head)
}

func exection(head *common.LinkList) {
	slow := head
	fast := head
	for {
		if fast == nil {
			break
		}
		if fast.Val != slow.Val {
			slow.Next = fast
			slow = slow.Next
		}
		fast = fast.Next
	}
	slow.Next = nil
	curr := head
	for {
		if curr == nil {
			break
		}
		fmt.Print(curr.Val, "")
		curr = curr.Next
	}
	fmt.Println("")
}
