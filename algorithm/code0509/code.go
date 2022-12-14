package code0509

import (
	"fmt"
)

func Exection() {
	n := 8
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	prev := 1
	curr := 1
	for i := 2; i <= n; i++ {
		sum := prev + curr
		prev = curr
		curr = sum
	}
	return curr
}
