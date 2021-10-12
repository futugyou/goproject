package main

import "fmt"

type tv struct {
	run bool
}

func (t *tv) on() {
	t.run = true
	fmt.Println("on")
}

func (t *tv) off() {
	t.run = false
	fmt.Println("off")
}
