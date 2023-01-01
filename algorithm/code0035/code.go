package code0035

import (
	"fmt"
)

func Exection() {
	nums := []int{1, 3, 5, 6}
	target := 7
	r := exection(nums, target)
	fmt.Println(r)
}

func exection(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := l + (r-l)/2
		switch {
		case target == nums[mid]:
			return mid
		case target < nums[mid]:
			r = mid - 1
		case target > nums[mid]:
			l = mid + 1
		}
	}

	return l
}
