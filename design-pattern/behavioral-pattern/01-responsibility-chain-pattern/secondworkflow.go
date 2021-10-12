package main

import "fmt"

type secondworkflow struct {
	execlevel int
	next      workflow
}

func (f *secondworkflow) exec(w work) {
	if f.execlevel >= w.worklevel {
		fmt.Println("secondworkflow exec the work")
	} else {
		if f.next == nil {
			fmt.Println("no workflow can exec the work")
		} else {
			f.next.exec(w)
		}
	}
}

func (f *secondworkflow) setnextflow(flow workflow) {
	f.next = flow
}
