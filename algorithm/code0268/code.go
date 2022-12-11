package code0268

import (
	"fmt"
)

func Exection() {
	nums := []int{0, 3, 1, 4}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) int {
	r := 0
	n := len(nums)
	r = r ^ n
	for i := 0; i < n; i++ {
		r = r ^ nums[i] ^ i
	}
	return r
}
