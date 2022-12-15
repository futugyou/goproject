package code0645

import (
	"fmt"
	"math"
)

func Exection() {
	nums := []int{1, 2, 2, 4}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) []int {
	due, miss := -1, -1
	for i := 0; i < len(nums); i++ {
		var index int = (int)(math.Abs(float64(nums[i])) - 1)
		if nums[index] < 0 {
			due = (int)(math.Abs(float64(nums[i])))
		} else {
			nums[i] *= -1
		}
	}
	for i := 0; i < len(nums); i++ {
		if nums[i] > 0 {
			miss = i + 1
			break
		}
	}
	return []int{due, miss}
}
