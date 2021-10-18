package main

import "fmt"

type istrategy interface {
	exec()
}

type strategyA struct {
}

func (s *strategyA) exec() {
	fmt.Println("strategyA")
}

type strategyB struct {
}

func (s *strategyB) exec() {
	fmt.Println("strategyB")
}

type strategyC struct {
}

func (s *strategyC) exec() {
	fmt.Println("strategyC")
}
