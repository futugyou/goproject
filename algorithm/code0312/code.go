package code0312

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	nums := []int{3, 1, 5, 8}
	r := exection(nums)
	fmt.Println(r)
}

func exection(nums []int) int {
	n := len(nums)
	points := make([]int, n+2)
	points[0], points[n+1] = 0, 0
	for i := 1; i <= n; i++ {
		points[i] = nums[i-1]
	}
	dp := make([][]int, n+2)
	for i := 0; i < n+2; i++ {
		dp[i] = make([]int, n+2)
	}
	for i := n; i >= 0; i-- {
		for j := i + 1; j < n+2; j++ {
			for k := i + 1; k < j; k++ {
				dp[i][j] = common.Max(dp[i][j], dp[i][k]+dp[k][j]+points[k]*points[i]*points[j])
			}
		}
	}
	return dp[0][n+1]
}
