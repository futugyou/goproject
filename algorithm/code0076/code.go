package code0076

import (
	"fmt"
	"math"
)

func Exection() {
	s := "abcabcbb"
	t := "bc"
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
	start := 0
	lenght := math.MaxInt
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
			if valid != len(need) {
				break
			}
			fmt.Println(left, right, need, dic)
			if right-left < lenght {
				start = left
				lenght = right - left
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
	if lenght == math.MaxInt {
		fmt.Println("")
	} else {
		s = string(r)
		fmt.Println(s[start : start+lenght])
	}
}
