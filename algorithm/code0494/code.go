package code0494

import "fmt"

func Exection() {
	nums := []int{1, 1, 1, 1, 1}
	target := 3
	exection(nums, target)
	fmt.Println(result)
}

func exection(nums []int, target int) {
	backtrack(nums, 0, target)
}

var result int = 0

func backtrack(nums []int, start, target int) {
	n := len(nums)
	if start == n {
		if target == 0 {
			result++
			return
		}
		return
	}
	target += nums[start]
	backtrack(nums, start+1, target)
	target -= nums[start]

	target -= nums[start]
	backtrack(nums, start+1, target)
	target += nums[start]
}
