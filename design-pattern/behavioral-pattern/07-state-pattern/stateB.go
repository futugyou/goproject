package main

import "fmt"

type stateB struct {
}

func (s *stateB) exec(con statecontext) {
	fmt.Println("state b")
	con.set(&stateC{})
	con.doexec()
}
