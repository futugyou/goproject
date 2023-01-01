package code0031

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{1, 2, 7, 4, 3, 1}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) []int {
	// check if is max array. eg 321
	ismax := true
	n := len(nums)
	for i := 0; i < n-1; i++ {
		if nums[i] < nums[i+1] {
			ismax = false
			break
		}
	}
	if ismax {
		common.Reverse(nums)
		return nums
	}
	var i int = 0
	for i = n - 1; i > 1; i-- {
		if nums[i-1] < nums[i] {
			i--
			break
		}
	}
	for j := n - 1; j > 0; j-- {
		if nums[j] > nums[i] {
			nums[j], nums[i] = nums[i], nums[j]
			break
		}
	}

	common.Reverse(nums[i+1:])
	return nums
}
