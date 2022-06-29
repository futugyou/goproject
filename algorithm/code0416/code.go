package code0416

import (
	"fmt"
)

func Exection() {
	nums := []int{1, 5, 11, 5}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) bool {
	sum := 0
	n := len(nums)
	for i := 0; i < n; i++ {
		sum += nums[i]
	}
	if sum%2 != 0 {
		return false
	}
	sum = sum / 2
	dp := make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]bool, sum+1)
		dp[i][0] = true
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= sum; j++ {
			if nums[i] > j {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i-1][j] || dp[i-1][j-nums[i]]
			}
		}
	}
	return dp[n][sum]
}
