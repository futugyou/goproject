package code0231

import (
	"fmt"
)

func Exection() {
	n := 99
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) bool {
	for {
		if n <= 0 {
			return false
		}
		if n&(n-1) == 0 {
			return true
		}
	}
}
