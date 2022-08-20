package code1201

import (
	"fmt"
	"math"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	n := 30
	a := 4
	b := 9
	c := 7
	r := exection(n, a, b, c)
	fmt.Println(r)
	r = exection2(n, a, b, c)
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

func exection2(n, a, b, c int) int {
	min := math.MinInt
	left := 0
	right := 2 * (int)(math.Pow10(9))
	for {
		if left > right {
			break
		}
		mid := left + (right-left)/2
		if f(mid, a, b, c) < n {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return min
}

// 容斥原理
func f(n, a, b, c int) int {
	seta := n / a
	setb := n / b
	setc := n / c
	setab := n / common.Lcm(a, b)
	setac := n / common.Lcm(a, c)
	setbc := n / common.Lcm(b, c)
	setabc := n / common.Lcm(common.Lcm(a, b), c)
	return seta + setb + setc - setab - setac - setbc + setabc
}
