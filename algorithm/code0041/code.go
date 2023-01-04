package code0041

import (
	"fmt"
)

func Exection() {
	nums := []int{7, 8, 9, 11, 12}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) int {
	n := len(nums)
	for i := 0; i < n; i++ {
		for nums[i] > 0 && nums[i] <= n && nums[i] != nums[nums[i]-1] {
			fmt.Println(i, nums)
			swap(&nums[i], &nums[nums[i]-1])
			fmt.Println(i, nums)
		}
	}
	fmt.Println(nums)
	for i := 0; i < n; i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}

	return n + 1
}

func swap(x, y *int) {
	*x, *y = *y, *x
}
