package main

import "fmt"

type stateA struct {
}

func (s *stateA) exec(con statecontext) {
	fmt.Println("state a")
	con.set(&stateB{})
	con.doexec()
}
