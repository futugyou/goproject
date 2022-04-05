package code0160

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	lista := common.NewLinkList()
	listb := common.NewLinkList()
	exection(lista, listb)
}

func exection(lista, listb common.LinkList) {
	a := &lista
	b := &listb
	for {
		if a == b {
			break
		}
		if a == nil {
			a = &listb
		} else {
			a = a.Next
		}
		if b == nil {
			b = &lista
		} else {
			b = b.Next
		}
	}
	if a != nil {
		a.Display()
	} else {
		fmt.Println("nil")
	}
}
