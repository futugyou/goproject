package code0167

import (
	"fmt"
)

func Exection() {
	nums := []int{2, 7, 11, 15}
	target := 9
	exection(nums, target)
}

func exection(nums []int, target int) {
	left := 0
	right := len(nums) - 1
	for {
		if left > right {
			break
		}
		sum := nums[left] + nums[right]
		if target == sum {
			fmt.Println(left, right)
			return
		} else if sum > target {
			right--
		} else if sum < target {
			left++
		}
	}
	fmt.Println(-1, -1)
}
