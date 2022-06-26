package code0309

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	prices := []int{1, 2, 3, 0, 2}
	r := exection(prices)
	fmt.Println(r)
}

func exection(prices []int) int {
	n := len(prices)
	dp := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, 2)
	}
	dp[0][0] = -prices[0]
	dp[0][1] = 0
	dp[1][0] = common.Max(dp[0][0], -prices[1])
	dp[1][1] = common.Max(dp[0][1], dp[1][0]+prices[1])

	for i := 2; i < n; i++ {
		dp[i][0] = common.Max(dp[i-1][0], dp[i-2][1]-prices[i])
		dp[i][1] = common.Max(dp[i-1][1], dp[i-1][0]+prices[i])
	}
	fmt.Println(dp)
	return dp[n-1][1]
}
