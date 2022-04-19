package code0022

import "fmt"

func Exection() {
	n := 3
	exection(n)
}

func exection(n int) {
	res := make([]string, 0)
	backtrack(n, n, "", &res)
	fmt.Println(res)
}

func backtrack(left, right int, s string, res *[]string) {
	if right < left || left < 0 || right < 0 {
		return
	}
	if left == 0 && right == 0 {
		*res = append(*res, s)
		return
	}
	s = s + "("
	backtrack(left-1, right, s, res)
	s = s[:len(s)-1]

	s = s + ")"
	backtrack(left, right-1, s, res)
	//s = s[:len(s)-1]
}
