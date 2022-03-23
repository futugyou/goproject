package code0344

import "fmt"

func Exection() {
	s := "[]int{2, 7, 11, 15}"

	exection(s)
}

func exection(s string) {
	r := []rune(s)
	left := 0
	right := len(r) - 1
	for {
		if left >= right {
			break
		}
		t := r[left]
		r[left] = r[right]
		r[right] = t
		left++
		right--
	}
	fmt.Println(string(r))
}
