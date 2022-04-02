package code0023

import (
	"github.com/futugyousuzu/goproject/algorithm/common"
	"github.com/gansidui/priority_queue"
)

func Exection() {
	nodes := make([]common.LinkList, 0)
	list1 := common.NewLinkList()
	list2 := common.NewLinkList()
	nodes = append(nodes, list1, list2)
	exection(nodes)
}

func exection(nodes []common.LinkList) {
	q := priority_queue.New()
	for _, v := range nodes {
		q.Push(v)
	}
	dummy := common.LinkList{Val: -1}
	p := &dummy
	for q.Len() > 0 {
		curr := q.Pop().(common.LinkList)
		p.Next = &curr
		if curr.Next != nil {
			q.Push(*curr.Next)
		}
		p = p.Next
	}
	r := dummy.Next
	r.Display()
}
