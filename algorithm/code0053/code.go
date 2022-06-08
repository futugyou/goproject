package code0053

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) int {
	n := len(nums)
	dp := make([]int, n)
	dp[0] = nums[0]
	result := 0
	for i := 1; i < n; i++ {
		dp[i] = common.Max(dp[i-1]+nums[i], nums[i])
	}
	for i := 0; i < n; i++ {
		result = common.Max(result, dp[i])
	}
	return result
}
