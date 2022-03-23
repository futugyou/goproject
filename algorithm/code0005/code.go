package code0005

import "fmt"

func Exection() {
	s := "babad"
	exection(s)
}

func exection(s string) {
	r := []rune(s)
	result := ""
	for i := 0; i < len(r); i++ {
		r1 := exec(r, i, i)
		r2 := exec(r, i, i+1)
		t := ""
		if len(r1) > len(r2) {
			t = string(r1)
		} else {
			t = string(r2)
		}
		if len(t) > len(result) {
			result = t
		}
	}
	fmt.Println(result)
}

func exec(s []rune, l, r int) []rune {
	for {
		if l < 0 || r >= len(s) || s[l] != s[r] {
			break
		}
		l--
		r++
	}
	l++
	return s[l:r]
}
