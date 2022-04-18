package code1541

import (
	"fmt"
)

func Exection() {
	s := "(()))"
	exection(s)
}

func exection(s string) {
	insert := 0
	rightneed := 0
	for _, v := range s {
		if string(v) == "(" {
			rightneed = rightneed + 2
			if rightneed%2 == 1 {
				insert++
				rightneed--
			}
		}
		if string(v) == ")" {
			rightneed--
			if rightneed < 0 {
				rightneed = 1
				insert++
			}
		}
	}
	fmt.Println(insert + rightneed)
}
