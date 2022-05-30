package code0698

import "fmt"

func Exection() {
	nums := []int{4, 3, 2, 3, 5, 2, 1}
	k := 4
	r := exection(nums, k)
	fmt.Println(r)
}

func exection(nums []int, k int) bool {
	if k > len(nums) {
		return false
	}
	sum := 0
	if sum%k != 0 {
		return false
	}
	used := make([]bool, len(nums))
	target := sum / k
	return backtrack(k, 0, nums, 0, used, target)
}

func backtrack(k, bucket int, nums []int, start int, used []bool, target int) bool {
	if k == 0 {
		return true
	}
	if bucket == target {
		return backtrack(k-1, 0, nums, 0, used, target)
	}
	for i := start; i < len(nums); i++ {
		if used[i] {
			continue
		}
		if nums[i]+bucket > target {
			continue
		}
		used[i] = true
		bucket += nums[i]
		if backtrack(k, bucket, nums, i+1, used, target) {
			return true
		}
		used[i] = false
		bucket -= nums[i]
	}
	return false
}
