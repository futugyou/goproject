package code0234

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	list := common.NewLinkList()
	r := exection(&list)
	fmt.Println(r)
}

func exection(list *common.LinkList) bool {
	fast := list
	slow := list
	for {
		if fast == nil || fast.Next == nil {
			break
		}
		fast = fast.Next.Next
		slow = slow.Next
	}
	if fast != nil {
		slow = slow.Next
	}
	left := list
	right := reverse(slow)
	for {
		if right == nil {
			break
		}
		if left.Val != right.Val {
			return false
		}
		left = left.Next
		right = right.Next
	}
	return true
}

func reverse(head *common.LinkList) *common.LinkList {
	if head == nil {
		return nil
	}
	if head.Next == nil {
		return head
	}
	newHdead := reverse(head.Next)
	head.Next.Next = head
	head.Next = nil
	return newHdead
}
