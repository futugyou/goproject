package code0026

import (
	"fmt"
	"sort"
)

func Exection() {
	nums := []int{1, 2, 3, 1, 2, 3, 4, 6, 7, 1, 2, 6}
	exection(nums)
}

func exection(nums []int) {
	sort.Ints(nums)
	slow := 0
	for fast := 0; fast < len(nums); fast++ {
		if nums[slow] != nums[fast] {
			slow++
			nums[slow] = nums[fast]
		}
	}
	fmt.Println(nums[0 : slow+1])
}
