package code0027

import (
	"fmt"
)

func Exection() {
	nums := []int{1, 2, 3, 1, 2, 3, 4, 6, 7, 1, 2, 6}
	val := 2
	exection(nums, val)
}

func exection(nums []int, val int) {
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[fast] != val {
			nums[slow] = nums[fast]
			slow++
		}
	}
	fmt.Println(nums[0:slow])
}
