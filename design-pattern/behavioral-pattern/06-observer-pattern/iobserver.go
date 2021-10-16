package main

import "fmt"

type iobserver interface {
	handler(sub isubject)
}

type observerA struct{}

func (o *observerA) handler(sub isubject) {
	v := sub.(*subject).getvalue()
	fmt.Println(v)
}
