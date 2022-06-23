package code0123

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	prices := []int{3, 3, 5, 0, 0, 3, 1, 4}
	r := exection(prices)
	fmt.Println(r)
}

func exection(prices []int) int {
	n := len(prices)
	dp := make([][][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([][]int, 3)
		for j := 0; j < 3; j++ {
			dp[i][j] = make([]int, 2)
		}
	}
	dp[0][0][0] = -prices[0]
	dp[0][0][1] = 0
	dp[0][1][0] = -prices[0]
	dp[0][1][1] = 0

	for i := 1; i < n; i++ {
		for j := 0; j < 2; j++ {
			if j == 0 {
				dp[i][j][0] = common.Max(dp[i-1][j][0], -prices[i])
				dp[i][j][1] = common.Max(dp[i-1][j][1], dp[i-1][j][0]+prices[i])
				continue
			}
			dp[i][j][0] = common.Max(dp[i-1][j][0], dp[i-1][j-1][1]-prices[i])
			dp[i][j][1] = common.Max(dp[i-1][j][1], dp[i-1][j][0]+prices[i])
		}
	}
	fmt.Println(dp)
	return dp[n-1][1][1]
}
