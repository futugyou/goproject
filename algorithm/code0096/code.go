package code0096

import (
	"fmt"
)

func Exection() {
	n := 3
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) int {
	if n <= 0 {
		return 0
	}
	return build(1, n)
}

func build(lo, hi int) int {
	res := 0
	if lo > hi {
		return 1
	}
	for i := lo; i <= hi; i++ {
		leftcount := build(lo, i-1)
		rightcunt := build(i+1, hi)
		res += leftcount * rightcunt
	}
	return res
}
