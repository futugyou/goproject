package code1201

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	n := 30
	a := 4
	b := 9
	c := 7
	r := exection(n, a, b, c)
	fmt.Println(r)
}

func exection(n, a, b, c int) int {
	min := 0
	vala := a
	valb := b
	valc := c
	p := 1
	for {
		if p > n {
			break
		}
		min = common.Min(vala, common.Min(valb, valc))
		p++
		if min == vala {
			vala += a
		}
		if min == valb {
			valb += b
		}
		if min == valc {
			valc += c
		}
	}
	return min
}
