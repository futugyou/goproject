package code0062

import (
	"fmt"
)

func Exection() {
	m := 7
	n := 3
	r := exection(m, n)
	fmt.Println(r)
}

func exection(m, n int) int {
	if m <= 0 || n <= 0 {
		return 0
	}
	if m <= 1 || n <= 1 {
		return 1
	}
	return exec(0, m-1, 0, n-1)
}

func exec(a, m, b, n int) int {
	if a >= m || b >= n {
		return 1
	}
	return exec(a+1, m, b, n) + exec(a, m, b+1, n)
}
