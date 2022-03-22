package code0567

import "fmt"

func Exection() {
	s1 := "ab"
	s2 := "eidbaooo"
	exection(s1, s2)
}

func exection(s1, s2 string) {
	need := make(map[rune]int)
	for _, v := range s1 {
		need[v]++
	}
	window := make(map[rune]int)
	left := 0
	right := 0
	result := false
	valid := 0
	s22 := []rune(s2)
	for {
		if right >= len(s22) || result {
			break
		}
		c := s22[right]
		right++
		if need[c] > 0 {
			window[c]++
			if need[c] == window[c] {
				valid++
			}
		}

		for {
			if right-left < len(s1) {
				break
			}
			if valid == len(need) {
				result = true
				break
			}
			d := s22[left]
			left++
			if need[d] > 0 {
				if need[d] == window[d] {
					valid--
				}
				window[d]--
			}
		}
	}
	fmt.Println(result)
}
