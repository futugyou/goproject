package code0283

import (
	"fmt"
)

func Exection() {
	nums := []int{1, 0, 3, 1, 0, 3, 4, 6, 7, 1, 0, 6}
	exection(nums)
}

func exection(nums []int) {
	val := 0
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != val {
			nums[slow] = nums[fast]
			slow++
		}
	}
	for i := slow; i < len(nums); i++ {
		nums[i] = 0
	}
	fmt.Println(nums)
}
