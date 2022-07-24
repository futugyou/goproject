package code0172

import (
	"fmt"
)

func Exection() {
	n := 126
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) int {
	res := 0
	d := 5
	for {
		if d > n {
			break
		}
		res += n / d
		d = d * 5
	}
	return res
}
