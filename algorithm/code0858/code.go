package code0858

import (
	"fmt"
)

func Exection() {
	p := 2
	q := 1
	r := exection(p, q)
	fmt.Println(r)
}

func exection(p int, q int) int {
	count := 1
	d := q
	for {
		if d%p == 0 {
			break
		}
		count++
		d += q
	}
	if count%2 == 0 {
		return 2
	}
	count = d / p
	if count%2 == 0 {
		return 0
	}
	return 1
}
