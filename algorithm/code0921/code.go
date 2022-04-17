package code0921

import (
	"fmt"
)

func Exection() {
	s := "())"
	exection(s)
}

func exection(s string) {
	left := 0
	right := 0
	for _, v := range s {
		if string(v) == "(" {
			right++
		}
		if string(v) == ")" {
			right--
			if right < 0 {
				right = 0
				left++
			}
		}
	}
	fmt.Println(left + right)
}
