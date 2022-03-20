package code0003

import "fmt"

func Exection() {
	s := "abcabcbb"
	exection(s)
}

func exection(s string) {
	r := []rune(s)
	left := 0
	right := 0
	maxlen := 0
	dic := make(map[rune]int)
	for {
		if right >= len(r) {
			break
		}
		c := r[right]
		right++
		dic[c]++
		for {
			if dic[c] <= 1 {
				break
			}
			d := r[left]
			left++
			dic[d]--
		}
		maxlen = max(maxlen, right-left)
	}
	fmt.Println(maxlen)
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}
