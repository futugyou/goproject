package code0191

import (
	"fmt"
)

func Exection() {
	n := 2
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) int {
	res := 0
	for {
		if n == 0 {
			break
		}
		n = n & (n - 1)
		res++
	}
	return res
}
