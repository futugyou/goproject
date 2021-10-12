package main

import "fmt"

type firstflow struct {
	execlevel int
	next      workflow
}

func (f *firstflow) exec(w work) {
	if f.execlevel >= w.worklevel {
		fmt.Println("firstflow exec the work")
	} else {
		if f.next == nil {
			fmt.Println("no workflow can exec the work")
		} else {
			f.next.exec(w)
		}
	}
}

func (f *firstflow) setnextflow(flow workflow) {
	f.next = flow
}
