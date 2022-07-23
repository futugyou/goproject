package code0136

import "fmt"

func Exection() {
	nums := []int{1, 1, 2, 2, 3, 4, 4, 5, 5}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) int {
	res := 0
	for i := 0; i < len(nums); i++ {
		res ^= nums[i]
	}
	return res
}
