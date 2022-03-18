package code0793

import (
	"fmt"
	"math"
)

func Exection() {
	k := 1
	exection(k)
}

func exection(k int) {
	left := leftbound(k)
	right := rightbound(k)
	fmt.Println(right, left, "=>", right-left+1)
}

func leftbound(k int) int {
	left := 0
	right := math.MaxInt
	for {
		if left > right {
			break
		}
		mid := left + (right-left)/2
		r := exec(mid)
		if r == k {
			right = mid - 1
		} else if r > k {
			right = mid - 1
		} else if r < k {
			left = mid + 1
		}
	}
	return left
}

func exec(mid int) int {
	r := 0
	for i := mid; i/5 > 0; i = i / 5 {
		r += i / 5
	}
	return r
}

func rightbound(k int) int {
	left := 0
	right := math.MaxInt
	for {
		if left > right {
			break
		}
		mid := left + (right-left)/2
		r := exec(mid)
		if r == k {
			left = mid + 1
		} else if r > k {
			right = mid - 1
		} else if r < k {
			left = mid + 1
		}
	}
	return right
}
