package code0096

import (
	"fmt"
)

var dic [][]int

func Exection() {
	n := 3
	dic = make([][]int, n+1)
	for i := 0; i < n+1; i++ {
		dic[i] = make([]int, n+1)
	}
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
	if dic[lo][hi] != 0 {
		return dic[lo][hi]
	}
	for i := lo; i <= hi; i++ {
		leftcount := build(lo, i-1)
		rightcunt := build(i+1, hi)
		res += leftcount * rightcunt
	}
	dic[lo][hi] = res
	return res
}
