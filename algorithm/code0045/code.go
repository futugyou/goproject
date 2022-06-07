package code0045

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{2, 3, 1, 1, 4}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) int {
	result := 0
	n := len(nums)
	fast := 0
	end := 0
	for i := 0; i < n-1; i++ {
		fast = common.Max(fast, nums[i]+i)
		if i == end {
			result++
			end = fast
		}
	}
	return result
}
