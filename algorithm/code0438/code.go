package code0438

import (
	"fmt"
)

func Exection() {
	s := "cbaebabacd"
	t := "abc"
	exection(s, t)
}

func exection(s, t string) {
	need := make(map[rune]int)
	for _, v := range t {
		need[v]++
	}
	r := []rune(s)
	left := 0
	right := 0
	valid := 0
	result := make([]int, 0)
	dic := make(map[rune]int)
	for {
		if right >= len(r) {
			break
		}
		c := r[right]
		right++

		if need[c] > 0 {
			dic[c]++
			if need[c] == dic[c] {
				valid++
			}
		}

		for {
			if right-left < len(t) {
				break
			}
			if valid == len(need) {
				result = append(result, left)
			}

			d := r[left]
			left++
			if need[d] > 0 {
				if need[d] == dic[d] {
					valid--
				}
				dic[d]--
			}
		}
	}
	fmt.Println(result)
}
