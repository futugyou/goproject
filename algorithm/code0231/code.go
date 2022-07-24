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
	if n <= 0 {
		return false
	}
	return n&(n-1) == 0
}
