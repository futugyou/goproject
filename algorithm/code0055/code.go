package code0055

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{2, 3, 1, 1, 4}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) bool {
	n := len(nums)
	fast := 0
	for i := 0; i < n-1; i++ {
		fast = common.Max(fast, nums[i]+i)
		if fast < i {
			return false
		}
	}
	return fast >= n-1
}
