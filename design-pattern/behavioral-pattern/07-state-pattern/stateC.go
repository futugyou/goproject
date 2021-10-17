package main

import "fmt"

type stateC struct {
}

func (s *stateC) exec(con statecontext) {
	fmt.Println("finish")
}
