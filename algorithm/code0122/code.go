package code0122

import (
	"fmt"

	"github.com/futugyousuzu/goproject/algorithm/common"
)

func Exection() {
	prices := []int{7, 1, 5, 3, 6, 4}
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
	for i := 1; i < n; i++ {
		dp[i][0] = common.Max(dp[i-1][0], dp[i-1][1]-prices[i])
		dp[i][1] = common.Max(dp[i-1][1], dp[i-1][0]+prices[i])
	}
	fmt.Println(dp)
	return dp[n-1][1]
}
