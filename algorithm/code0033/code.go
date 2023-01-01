package code0033

import (
	"fmt"
)

func Exection() {
	nums := []int{4, 5, 6, 7, 0, 1, 2}
	target := 5
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
		case nums[l] <= nums[mid]:
			if nums[l] <= target && target < nums[mid] {
				r = mid - 1
			} else {
				l = mid + 1
			}
		case nums[mid] <= nums[r]:
			if nums[mid] < target && target <= nums[r] {
				l = mid + 1
			} else {
				r = mid - 1
			}
		}
	}

	return -1
}
