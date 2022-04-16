package code0020

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	s := "()[]{}"
	exection(s)
}

func exection(s string) {
	stack := common.NewStack()
	for _, v := range s {
		if string(v) == "{" || string(v) == "(" || string(v) == "[" {
			stack.Push(string(v))
		} else {
			if !stack.Empty() && leftof(string(v)) == stack.Peak() {
				stack.Pop()
			} else {
				fmt.Println(false)
				return
			}
		}
	}
	fmt.Println(stack.Empty())
}

func leftof(s string) string {
	if s == "}" {
		return "{"
	}
	if s == "]" {
		return "["
	}
	return "("
}
