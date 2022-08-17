package code0263

import (
	"fmt"
)

func Exection() {
	n := 30
	r := exection(n)
	fmt.Println(r)
}

func exection(n int) bool {
	if n <= 0 {
		return false
	}
	for {
		if n%2 != 0 {
			break
		}
		n = n / 2
	}
	for {
		if n%3 != 0 {
			break
		}
		n = n / 3
	}
	for {
		if n%5 != 0 {
			break
		}
		n = n / 5
	}
	return n == 1
}
